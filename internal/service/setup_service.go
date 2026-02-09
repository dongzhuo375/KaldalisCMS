package service

import (
	"KaldalisCMS/internal/infra/model"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
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

func (s *SetupService) Install(cfg SetupConfig) error {
	// 1. Try Connect
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	// 2. Migrate
	if err := db.AutoMigrate(&model.User{}, &model.Category{}, &model.Tag{}, &model.Post{}, &model.SystemSetting{}); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	// 3. Create Admin User
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPass), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
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
			return fmt.Errorf("failed to create admin user: %w", err)
		}
	}

	// 4. Mark System as Installed
	now := time.Now()
	setting := model.SystemSetting{
		Installed:   true,
		SiteName:    cfg.SiteName,
		InstalledAt: &now,
	}
	if err := db.FirstOrCreate(&setting, model.SystemSetting{ID: 1}).Error; err != nil {
		return fmt.Errorf("failed to set system setting: %w", err)
	}
	
	if err := db.Model(&setting).Updates(model.SystemSetting{SiteName: cfg.SiteName, Installed: true}).Error; err != nil {
		return fmt.Errorf("failed to update system setting: %w", err)
	}

	// 5. Save Config to file (Using basic types in callback)
	if s.SaveConfigFunc != nil {
		if err := s.SaveConfigFunc(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName); err != nil {
			return fmt.Errorf("failed to save configuration file: %w", err)
		}
	}

	// 6. Trigger Hot Reload
	if s.ReloadFunc != nil {
		go func() {
			time.Sleep(1 * time.Second)
			s.ReloadFunc()
		}()
	}

	return nil
}