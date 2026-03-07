package service

import (
	"KaldalisCMS/internal/infra/model"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
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

	// 权限配置标志
	AllowAnonymousRead bool
	AdminFullAccess    bool
}

type SetupService struct {
	SaveConfigFunc func(host string, port int, user, pass, dbname string) error
	ReloadFunc     func() error
	Enforcer       *casbin.Enforcer // 用于初始化权限
}

func NewSetupService(save func(string, int, string, string, string) error, reload func() error) *SetupService {
	return &SetupService{SaveConfigFunc: save, ReloadFunc: reload}
}

// SetEnforcer 允许在安装开始前动态注入 Enforcer
func (s *SetupService) SetEnforcer(enforcer *casbin.Enforcer) {
	s.Enforcer = enforcer
}

var reIdentifier = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// ValidateDatabase 采用多级探测机制，确保即便默认管理库缺失也能正常初始化。
func (s *SetupService) ValidateDatabase(host string, port int, user, pass, dbname string) error {
	if dbname == "" {
		return fmt.Errorf("数据库名称不能为空")
	}

	// --- 1. 尝试直接连接目标库 (第一优先级) ---
	targetDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		host, port, user, pass, dbname)
	
	fmt.Printf("[SETUP] 正在验证目标数据库: %s\n", dbname)
	db, err := gorm.Open(postgres.Open(targetDSN), &gorm.Config{})
	if err == nil {
		if sqlDB, err := db.DB(); err == nil {
			err = sqlDB.Ping()
			sqlDB.Close()
			if err == nil {
				fmt.Printf("[SETUP] 目标数据库 [%s] 已存在且可连接，跳过创建流程。\n", dbname)
				return nil 
			}
		}
	}

	// --- 2. 目标库不存在，尝试通过管理库创建 ---
	adminDBs := []string{"postgres", "template1"}
	var adminDB *gorm.DB
	var lastErr error

	for _, adminName := range adminDBs {
		fmt.Printf("[SETUP] 尝试通过管理库 [%s] 进行自动创建...\n", adminName)
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			host, port, user, pass, adminName)
		
		adminDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		lastErr = err
	}

	if adminDB == nil {
		return fmt.Errorf("无法连接到任何管理数据库 (postgres/template1): %w", lastErr)
	}
	// 确保无论如何管理库连接都会关闭
	defer func() {
		if sqlDB, err := adminDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	if !reIdentifier.MatchString(dbname) {
		return fmt.Errorf("无效的数据库名: %s (只能包含字母、数字和下划线)", dbname)
	}

	var exists int
	adminDB.Raw("SELECT 1 FROM pg_database WHERE datname = ?", dbname).Scan(&exists)
	if exists == 0 {
		fmt.Printf("[SETUP] 数据库 [%s] 不存在，执行创建指令...\n", dbname)
		if err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname)).Error; err != nil {
			return fmt.Errorf("尝试自动创建数据库失败: %w", err)
		}
	}

	// --- 3. 最后一次终极验证 ---
	finalDB, err := gorm.Open(postgres.Open(targetDSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库创建成功但无法连接: %w", err)
	}
	if sqlDB, err := finalDB.DB(); err == nil {
		defer sqlDB.Close()
		return sqlDB.Ping()
	}

	return fmt.Errorf("数据库验证逻辑发生未知错误")
}

func (s *SetupService) Install(cfg SetupConfig) error {
	// 执行多级预检
	if err := s.ValidateDatabase(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName); err != nil {
		return err
	}

	// 正式连接
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// 迁移表结构
	if err := db.AutoMigrate(&model.User{}, &model.Category{}, &model.Tag{}, &model.Post{}, &model.SystemSetting{}, &model.MediaAsset{}, &model.PostAsset{}); err != nil {
		return fmt.Errorf("Schema 迁移失败: %w", err)
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
			return fmt.Errorf("管理员创建失败: %w", err)
		}
	}

	// 保存站点设置
	now := time.Now()
	setting := model.SystemSetting{
		ID:          1,
		Installed:   true,
		SiteName:    cfg.SiteName,
		InstalledAt: &now,
	}
	if err := db.Save(&setting).Error; err != nil {
		return fmt.Errorf("系统设置持久化失败: %w", err)
	}

	// --- 权限初始化 (Casbin) ---
	// 初始化 GORM 适配器并加载模型配置
	adapter, _ := gormadapter.NewAdapterByDB(db)
	enforcer, err := casbin.NewEnforcer("cmd/configs/casbin_model.conf", adapter)
	if err == nil {
		enforcer.EnableAutoSave(true)

		// 1. 如果开启了管理员全能权限，为 super_admin 分配通配符策略
		if cfg.AdminFullAccess {
			enforcer.AddPolicy("super_admin", "/api/v1/*", "*")
			log.Printf("[SETUP] 已为超级管理员角色分配全量 API 权限")
		}

		// 2. 如果开启了访客阅读，分配基础 GET 权限
		if cfg.AllowAnonymousRead {
			enforcer.AddPolicy("anonymous", "/api/v1/posts", "GET")
			enforcer.AddPolicy("anonymous", "/api/v1/posts/:id", "GET")
			log.Printf("[SETUP] 已开启访客匿名阅读权限策略")
		}

		// 3. 将当前创建的管理员账号绑定到 super_admin 角色
		enforcer.AddGroupingPolicy(cfg.AdminUser, "super_admin")

		log.Printf("[SETUP] 初始 Casbin 策略处理完成 (Admin: %s -> super_admin)", cfg.AdminUser)
	} else {
		log.Printf("[WARN] 权限策略初始化失败 (可能是配置文件缺失): %v", err)
	}

	// 保存配置文件
	if s.SaveConfigFunc != nil {
		if err := s.SaveConfigFunc(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName); err != nil {
			return fmt.Errorf("YAML 配置文件更新失败: %w", err)
		}
	}

	// 热重载
	if s.ReloadFunc != nil {
		go func() {
			time.Sleep(1 * time.Second)
			s.ReloadFunc()
		}()
	}

	return nil
}
