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

	// --- 1. 移除数据库默认值 ---
	// 为了让系统在没有配置时能正确进入 SETUP MODE，我们不再为数据库设置默认值。
	// 只有 SSLMode 和 TimeZone 这种可选参数保留默认。
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.timezone", "Asia/Shanghai")

	// Media 默认值保持
	v.SetDefault("media.upload_dir", filepath.FromSlash("./data/uploads"))
	v.SetDefault("media.max_upload_size_mb", int64(50))
	v.SetDefault("media.max_filename_bytes", 180)

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("未找到配置文件，将使用环境或等待初始化安装。")
		} else {
			log.Fatalf("读取配置文件出错: %s \n", err)
		}
	}

	// Environment variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Unmarshal into struct
	if err := v.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("无法解析配置到结构体: %v", err)
	}

	// 进一步细化 Auth 配置（如 JWT Secret 等）
	refinedAuth, err := auth.LoadConfig(v)
	if err != nil {
		log.Printf("Auth 配置初始化警告 (可能尚未安装): %v", err)
	} else {
		AppConfig.Auth = *refinedAuth
	}

	log.Println("配置加载流程结束。")
}

func GetDatabaseDSN() string {
	db := AppConfig.Database
	// 如果关键字段为空，返回空 DSN 或报错，让连接失败触发 SETUP MODE
	if db.Host == "" || db.User == "" || db.DBName == "" {
		return ""
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode, db.TimeZone)
}

// SaveDatabaseConfig 仅更新数据库部分，保留现有的 jwt, media 等配置
func SaveDatabaseConfig(host string, port int, user, pass, dbname string) error {
	v := viper.GetViper()

	// 1. 更新数据库相关字段
	v.Set("database.host", host)
	v.Set("database.port", port)
	v.Set("database.user", user)
	v.Set("database.password", pass)
	v.Set("database.dbname", dbname)

	// 更新内存中的 AppConfig，确保热重启时拿的是最新的
	AppConfig.Database.Host = host
	AppConfig.Database.Port = port
	AppConfig.Database.User = user
	AppConfig.Database.Password = pass
	AppConfig.Database.DBName = dbname

	// 2. 确保配置目录存在
	configPath := "./cmd/configs/config.yaml"
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("创建配置目录失败: %w", err)
		}
	}

	// 3. 写入文件
	// 使用 WriteConfigAs 而不是 WriteConfig，因为初始安装时可能文件根本不存在。
	// Viper 会把当前内存里已有的所有信息（包括之前读取到的 jwt, media）一并写回。
	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	log.Printf("配置已成功更新至 %s", configPath)
	return nil
}
