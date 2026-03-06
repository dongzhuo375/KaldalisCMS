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

// reIdentifier 用于基础的标识符校验，防止 SQL 注入
var reIdentifier = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func (s *SetupService) Install(cfg SetupConfig) error {
	// --- 1. 数据库自动拨备 (Provisioning) ---
	// 先连到默认的 postgres 管理库，检查并创建目标数据库
	adminDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass)

	adminDB, err := gorm.Open(postgres.Open(adminDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("无法连接到 Postgres 管理库 (请检查账号密码或主机): %w", err)
	}

	// 确保目标数据库名合法
	if !reIdentifier.MatchString(cfg.DbName) {
		return fmt.Errorf("无效的数据库名: %s (仅限字母数字下划线)", cfg.DbName)
	}

	// 检查库是否存在
	var exists int
	adminDB.Raw("SELECT 1 FROM pg_database WHERE datname = ?", cfg.DbName).Scan(&exists)
	if exists == 0 {
		// CREATE DATABASE 不支持参数化查询，所以必须用上面的正则校验过
		if err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DbName)).Error; err != nil {
			return fmt.Errorf("自动创建数据库失败: %w", err)
		}
	}

	// 关闭管理库连接
	if sqlDB, err := adminDB.DB(); err == nil {
		sqlDB.Close()
	}

	// --- 2. 正式连接业务库 ---
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接业务数据库失败: %w", err)
	}

	// 严谨探测连接是否真的可用
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层连接池失败: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库 Ping 探测失败: %w", err)
	}

	// --- 3. 数据库迁移 (Migrate) ---
	if err := db.AutoMigrate(&model.User{}, &model.Category{}, &model.Tag{}, &model.Post{}, &model.SystemSetting{}); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	// --- 4. 创建超级管理员 ---
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPass), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
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
			return fmt.Errorf("创建管理员账号失败: %w", err)
		}
	}

	// --- 5. 写入安装标记与站点名 ---
	now := time.Now()
	setting := model.SystemSetting{
		ID:          1,
		Installed:   true,
		SiteName:    cfg.SiteName,
		InstalledAt: &now,
	}
	// 使用 Save 确保根据 ID 1 更新或创建记录
	if err := db.Save(&setting).Error; err != nil {
		return fmt.Errorf("更新系统设置失败: %w", err)
	}

	// --- 6. 持久化配置文件 ---
	if s.SaveConfigFunc != nil {
		if err := s.SaveConfigFunc(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName); err != nil {
			return fmt.Errorf("持久化配置文件失败: %w", err)
		}
	}

	// --- 7. 触发热重启 ---
	if s.ReloadFunc != nil {
		go func() {
			// 给响应返回留出一点时间
			time.Sleep(1 * time.Second)
			s.ReloadFunc()
		}()
	}

	return nil
}
