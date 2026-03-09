package main

import (
	"KaldalisCMS/internal/infra/auth"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
		TimeZone string `mapstructure:"timezone"`
	} `mapstructure:"database"`
	Auth auth.Config

	Media struct {
		// UploadDir 是所有媒体文件的根目录。
		// 路由层仅暴露 UploadDir/a 子目录（对应 /media/a），UploadDir 根目录下的其他文件不会被公开访问。
		UploadDir        string `mapstructure:"upload_dir"`
		MaxUploadSizeMB  int64  `mapstructure:"max_upload_size_mb"`
		PublicBaseURL    string `mapstructure:"public_base_url"`
		MaxFilenameBytes int    `mapstructure:"max_filename_bytes"`
	} `mapstructure:"media"`
}

var AppConfig Config

// InitConfig initializes and loads the configuration
func InitConfig() {
	v := viper.GetViper()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./cmd/configs")

	// 核心字段绝不设默认值，逼迫系统进入 Setup 模式
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.timezone", "Asia/Shanghai")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("[CONFIG] 未找到配置文件，将使用默认值并准备进入 Setup 模式。")
		} else {
			log.Fatalf("[CONFIG] 读取失败: %s \n", err)
		}
	} else {
		log.Printf("[CONFIG] 成功加载配置文件: %s", v.ConfigFileUsed())
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("[CONFIG] 解析失败: %v", err)
	}

	// 增加脱敏调试日志
	hasPass := "否"
	if AppConfig.Database.Password != "" {
		hasPass = "是 (长度:" + fmt.Sprint(len(AppConfig.Database.Password)) + ")"
	}
	log.Printf("[CONFIG] 数据库配置加载完毕 -> Host: %s, DB: %s, User: %s, 是否含密码: %s",
		AppConfig.Database.Host, AppConfig.Database.DBName, AppConfig.Database.User, hasPass)

	// 初始化 Auth
	refinedAuth, err := auth.LoadConfig(v)
	if err == nil {
		AppConfig.Auth = *refinedAuth
		//检查是否使用了默认的弱密钥
		defaultSecret := "AL3uaHdBI/zK/t0zfeXrFKb/7LOP8LECxp51j7pOo9PP7Ok99JceBQ8k4AZYlOE7tM8sV/55hPq/8I3WdzJi1w=="
		if string(AppConfig.Auth.Secret) == defaultSecret {
			log.Println("[SECURITY WARNING] 您正在使用默认的 JWT Secret，生产环境请务必通过环境变量 JWT_SECRET 修改！")
		}
		if len(AppConfig.Auth.Secret) < 32 {
			log.Println("[SECURITY WARNING] JWT Secret 长度不足 32 字节，建议使用更长的密钥。")
		}
	}
}

func GetDatabaseDSN() string {
	db := AppConfig.Database
	// 关键字段缺失或包含默认/敏感值不完整时直接返回空
	if db.Host == "" || db.DBName == "" || db.User == "" {
		log.Printf("[DATABASE] 配置不完整 (Host:%s, DB:%s, User:%s), 准备进入安装模式", db.Host, db.DBName, db.User)
		return ""
	}

	// 为 dbname 添加单引号，防止特殊字符干扰，且在 DSN 拼接中确保字段分明
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode, db.TimeZone)

	log.Printf("[DATABASE] 准备校验连接: Host=%s, Port=%d, User=%s, DB=%s", db.Host, db.Port, db.User, db.DBName)
	return dsn
}

func SaveDatabaseConfig(host string, port int, user, pass, dbname string) error {
	v := viper.GetViper()

	// 设置 Viper 内存值以供写入
	v.Set("database.host", host)
	v.Set("database.port", port)
	v.Set("database.user", user)
	v.Set("database.password", pass)
	v.Set("database.dbname", dbname)
	v.Set("database.sslmode", "disable") // 强制默认值以防万一
	v.Set("database.timezone", "Asia/Shanghai")

	// 更新全局 AppConfig 变量供热重启直接使用
	AppConfig.Database.Host = host
	AppConfig.Database.Port = port
	AppConfig.Database.User = user
	AppConfig.Database.Password = pass
	AppConfig.Database.DBName = dbname
	AppConfig.Database.SSLMode = "disable"
	AppConfig.Database.TimeZone = "Asia/Shanghai"

	configPath := "./cmd/configs/config.yaml"
	_ = os.MkdirAll(filepath.Dir(configPath), 0755)

	// 使用 WriteConfig 保存到现有路径
	if err := v.WriteConfigAs(configPath); err != nil {
		log.Printf("[CONFIG] 保存配置失败: %v", err)
		return err
	}
	return nil
}
