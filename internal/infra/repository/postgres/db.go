package repository

import (
	model2 "KaldalisCMS/internal/infra/model"
	"errors"
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, errors.New("database DSN is empty")
	}

	// 尝试从 DSN 中提取预期的数据库名用于校验
	expectedDB := ""
	parts := strings.Split(dsn, " ")
	for _, p := range parts {
		if strings.HasPrefix(p, "dbname=") {
			expectedDB = strings.TrimPrefix(p, "dbname=")
			break
		}
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// --- 核心校验：确保连上的库就是我们要的那个 ---
	var currentDB string
	db.Raw("SELECT current_database()").Scan(&currentDB)
	
	if expectedDB != "" && currentDB != expectedDB {
		// 彻底关闭这个错误的连接
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
		return nil, fmt.Errorf("安全拦截：预期连接数据库 [%s], 但实际连上了 [%s]。请检查配置或清理 Postgres 环境变量。", expectedDB, currentDB)
	}

	log.Printf("[DATABASE] 成功连接并校验数据库: %s", currentDB)
	// --- END 核心校验 ---

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&model2.User{},
		&model2.Category{},
		&model2.Tag{},
		&model2.Post{},
		&model2.SystemSetting{},
		&model2.MediaAsset{},
		&model2.PostAsset{},
	)
	if err != nil {
		log.Printf("Failed to auto-migrate database: %v", err)
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}
	fmt.Println("Database schema migrated successfully.")

	return db, nil
}
