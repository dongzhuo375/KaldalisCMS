package service

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/infra/model"
	"fmt"
	"log"
	"regexp"
	"strings"
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

	// 细粒度权限配置标志
	AllowAnonymousRead bool
	AdminFullAccess    bool
	AdminCanDelete     bool
	UserCanUpload      bool
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

func validateInstallConfig(cfg SetupConfig) error {
	if strings.TrimSpace(cfg.DbHost) == "" || cfg.DbPort <= 0 || strings.TrimSpace(cfg.DbUser) == "" || strings.TrimSpace(cfg.DbName) == "" {
		return fmt.Errorf("%w: invalid database setup parameters", core.ErrInvalidInput)
	}
	if strings.TrimSpace(cfg.SiteName) == "" || strings.TrimSpace(cfg.AdminUser) == "" || strings.TrimSpace(cfg.AdminEmail) == "" || cfg.AdminPass == "" {
		return fmt.Errorf("%w: missing required setup fields", core.ErrInvalidInput)
	}
	if len([]byte(cfg.AdminPass)) > 72 {
		return fmt.Errorf("%w: admin password exceeds bcrypt limit", core.ErrInvalidInput)
	}
	return nil
}

// ValidateDatabase 采用多级探测机制，确保即便默认管理库缺失也能正常初始化。
func (s *SetupService) ValidateDatabase(host string, port int, user, pass, dbname string) error {
	if strings.TrimSpace(dbname) == "" {
		return normalizeServiceErrorWithOpMsg("setup.validate_database.input", "database name is required", core.ErrInvalidInput)
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
		return normalizeServiceErrorWithOpMsg("setup.validate_database.admin_connect", "connect to admin database failed", fmt.Errorf("%w: %v", core.ErrInternalError, lastErr))
	}
	// 确保无论如何管理库连接都会关闭
	defer func() {
		if sqlDB, err := adminDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	if !reIdentifier.MatchString(dbname) {
		return normalizeServiceErrorWithOpMsg("setup.validate_database.identifier", "database name format is invalid", core.ErrInvalidInput)
	}

	var exists int
	adminDB.Raw("SELECT 1 FROM pg_database WHERE datname = ?", dbname).Scan(&exists)
	if exists == 0 {
		fmt.Printf("[SETUP] 数据库 [%s] 不存在，执行创建指令...\n", dbname)
		if err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname)).Error; err != nil {
			return normalizeServiceErrorWithOpMsg("setup.validate_database.create_db", "auto create database failed", err)
		}
	}

	// --- 3. 最后一次终极验证 ---
	finalDB, err := gorm.Open(postgres.Open(targetDSN), &gorm.Config{})
	if err != nil {
		return normalizeServiceErrorWithOpMsg("setup.validate_database.final_open", "open target database after creation failed", err)
	}
	if sqlDB, err := finalDB.DB(); err == nil {
		defer sqlDB.Close()
		if err := sqlDB.Ping(); err != nil {
			return normalizeServiceErrorWithOpMsg("setup.validate_database.final_ping", "ping target database failed", err)
		}
		return nil
	}

	return normalizeServiceErrorWithOpMsg("setup.validate_database.final_state", "database validation ended in unknown state", core.ErrInternalError)
}

func (s *SetupService) Install(cfg SetupConfig) error {
	if err := validateInstallConfig(cfg); err != nil {
		return normalizeServiceErrorWithOpMsg("setup.install.validate_config", "validate setup config failed", err)
	}

	// 执行多级预检
	if err := s.ValidateDatabase(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName); err != nil {
		return normalizeServiceErrorWithOpMsg("setup.install.validate_database", "validate database before install failed", err)
	}

	// 正式连接
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return normalizeServiceErrorWithOpMsg("setup.install.open_db", "open install target database failed", err)
	}

	// 迁移表结构
	if err := db.AutoMigrate(&model.User{}, &model.Category{}, &model.Tag{}, &model.Post{}, &model.SystemSetting{}, &model.MediaAsset{}, &model.PostAsset{}); err != nil {
		return normalizeServiceErrorWithOpMsg("setup.install.migrate", "schema migration failed", err)
	}

	// 创建管理员
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPass), bcrypt.DefaultCost)
	if err != nil {
		return normalizeServiceErrorWithOpMsg("setup.install.hash_admin_password", "hash admin password failed", core.ErrInvalidInput)
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
			return normalizeServiceErrorWithOpMsg("setup.install.create_admin", "create admin user failed", err)
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
		return normalizeServiceErrorWithOpMsg("setup.install.save_setting", "persist system setting failed", err)
	}

	// --- 权限初始化 (Casbin RBAC 细粒度模板) ---
	adapter, _ := gormadapter.NewAdapterByDB(db)
	enforcer, err := casbin.NewEnforcer("cmd/configs/casbin_model.conf", adapter)
	if err == nil {
		enforcer.EnableAutoSave(true)

		// super_admin is handled entirely by the matcher hardcode (r.sub == "super_admin"),
		// so no explicit policy is needed.

		// 1. [Role: admin] - 内容管理员
		// 后台文章管理统一走 /api/v1/admin/posts，公共 /posts 只承担已发布内容分发。
		adminRules := [][]string{
			{"admin", "/api/v1/admin/posts", "GET"},
			{"admin", "/api/v1/admin/posts", "POST"},
			{"admin", "/api/v1/admin/posts/:id", "GET"},
			{"admin", "/api/v1/admin/posts/:id", "PUT"},
			{"admin", "/api/v1/admin/posts/:id/publish", "POST"},
			{"admin", "/api/v1/admin/posts/:id/draft", "POST"},
			{"admin", "post", "list:any"},
			{"admin", "post", "read:any"},
			{"admin", "post", "update:any"},
			{"admin", "post", "publish"},
			{"admin", "post", "unpublish"},
			{"admin", "post", "delete"},
			{"admin", "/api/v1/media", "POST"},
			{"admin", "/api/v1/tags", "POST"},
			{"admin", "/api/v1/tags/:id", "PUT"},
			{"admin", "/api/v1/categories", "POST"},
			{"admin", "/api/v1/categories/:id", "PUT"},
		}
		//理想状态为初始化建站时载入的权限策略，实际上在router存在时就会自动加载到内存（我看未必），所以这里的AddPolicy更多是为了确保安装流程的完整性和可预见性。
		enforcer.AddPolicies(adminRules)

		if cfg.AdminCanDelete {
			enforcer.AddPolicy("admin", "/api/v1/admin/posts/:id", "DELETE")
			enforcer.AddPolicy("admin", "/api/v1/media/:id", "DELETE")
			enforcer.AddPolicy("admin", "/api/v1/tags/:id", "DELETE")
			enforcer.AddPolicy("admin", "/api/v1/categories/:id", "DELETE")
		}

		// 3. [Role: user] - 普通注册用户
		// 用户可进入文章后台管理自己的草稿，但最终的数据范围仍由服务层限制为“仅本人 Draft”。
		userRules := [][]string{
			{"user", "/api/v1/posts", "GET"},
			{"user", "/api/v1/posts/:id", "GET"},
			{"user", "/api/v1/admin/posts", "GET"},
			{"user", "/api/v1/admin/posts", "POST"},
			{"user", "/api/v1/admin/posts/:id", "GET"},
			{"user", "/api/v1/admin/posts/:id", "PUT"},
			{"user", "post:draft", "create"},
			{"user", "post:draft", "list:own"},
			{"user", "post:draft", "read:own"},
			{"user", "post:draft", "update:own"},
			{"user", "/api/v1/media", "GET"},
		}
		enforcer.AddPolicies(userRules)

		if cfg.UserCanUpload {
			enforcer.AddPolicy("user", "/api/v1/media", "POST")
		}

		// 4. [Role: anonymous] - 匿名访客
		if cfg.AllowAnonymousRead {
			enforcer.AddPolicy("anonymous", "/api/v1/posts", "GET")
			enforcer.AddPolicy("anonymous", "/api/v1/posts/:id", "GET")
		}

		// 4. [Inheritance] - 角色继承
		enforcer.AddGroupingPolicy("admin", "user")
		enforcer.AddGroupingPolicy("super_admin", "admin")

		// NOTE: No need to bind username → super_admin via AddGroupingPolicy.
		// The Casbin subject is always the role string from JWT claims (e.g. "super_admin"),
		// not the username. The admin user already has Role="super_admin" in the users table.

		log.Printf("[SETUP] 细粒度 RBAC 体系注入完成 (AdminDelete:%v, UserUpload:%v)", cfg.AdminCanDelete, cfg.UserCanUpload)
	} else {
		log.Printf("[WARN] 权限策略初始化失败: %v", err)
	}

	// 保存配置文件
	if s.SaveConfigFunc != nil {
		if err := s.SaveConfigFunc(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName); err != nil {
			return normalizeServiceErrorWithOpMsg("setup.install.save_config", "persist setup config failed", err)
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
