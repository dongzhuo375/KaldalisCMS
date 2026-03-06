package service

import (
	"KaldalisCMS/internal/infra/model"
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SetupConfig struct {
	DbHost     string
	DbPort     int
	DbUser     string
	DbPass     string
	DbName     string
	SiteName   string
	AdminUser  string
	AdminPass  string
	AdminEmail string
}

type SetupService struct {
	SaveConfigFunc func(host string, port int, user, pass, dbname string) error
	ReloadFunc     func() error
}

func NewSetupService(save func(string, int, string, string, string) error, reload func() error) *SetupService {
	return &SetupService{SaveConfigFunc: save, ReloadFunc: reload}
}

var reIdentifier = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// ValidateDatabase 仅执行：连管理库 -> 检查并建库 -> 连目标库 -> Ping 探测。
// 它不保存配置，也不重启系统，纯粹用于预检。
func (s *SetupService) ValidateDatabase(host string, port int, user, pass, dbname string) error {
	// 1. 连管理库
	adminDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable TimeZone=Asia/Shanghai",
		host, port, user, pass)

	adminDB, err := gorm.Open(postgres.Open(adminDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("无法连接到数据库服务 (请检查账号密码或主机): %w", err)
	}

	// 2. 校验库名
	if !reIdentifier.MatchString(dbname) {
		return fmt.Errorf("无效的数据库名: %s (仅限字母数字下划线)", dbname)
	}

	// 3. 检查并自动建库
	var exists int
	adminDB.Raw("SELECT 1 FROM pg_database WHERE datname = ?", dbname).Scan(&exists)
	if exists == 0 {
		if err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname)).Error; err != nil {
			return fmt.Errorf("自动创建数据库失败: %w", err)
		}
	}

	// 关闭管理库连接
	if sqlDB, err := adminDB.DB(); err == nil {
		sqlDB.Close()
	}

	// 4. 正式连目标库探活
	targetDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		host, port, user, pass, dbname)

	db, err := gorm.Open(postgres.Open(targetDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("目标数据库连接失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层连接池失败: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库 Ping 探测失败: %w", err)
	}
	sqlDB.Close()

	return nil
}

func (s *SetupService) Install(cfg SetupConfig) error {
	// 复用预检逻辑
	if err := s.ValidateDatabase(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName); err != nil {
		return err
	}

	// 连接并开始安装
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// 数据库迁移
	if err := db.AutoMigrate(&model.User{}, &model.Category{}, &model.Tag{}, &model.Post{}, &model.SystemSetting{}, &model.MediaAsset{}, &model.PostAsset{}); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 创建管理员
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := model.User{
		Username:  cfg.AdminUser,
		Email:     cfg.AdminEmail,
		Password:  string(hashedPassword),
		Role:      "super_admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var count int64
	db.Model(&model.User{}).Where("username = ?", admin.Username).Count(&count)
	if count == 0 {
		if err := db.Create(&admin).Error; err != nil {
			return fmt.Errorf("创建管理员失败: %w", err)
		}
	}

	// 标记安装成功
	now := time.Now()
	setting := model.SystemSetting{
		ID:          1,
		Installed:   true,
		SiteName:    cfg.SiteName,
		InstalledAt: &now,
	}
	if err := db.Save(&setting).Error; err != nil {
		return fmt.Errorf("保存系统设置失败: %w", err)
	}

	// 保存配置
	if s.SaveConfigFunc != nil {
		if err := s.SaveConfigFunc(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName); err != nil {
			return fmt.Errorf("持久化配置失败: %w", err)
		}
	}

	// 触发重启
	if s.ReloadFunc != nil {
		go func() {
			time.Sleep(1 * time.Second)
			s.ReloadFunc()
		}()
	}

	return nil
}
