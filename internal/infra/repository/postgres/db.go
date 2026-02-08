package repository

import (
	model2 "KaldalisCMS/internal/infra/model"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Database connection established successfully.")

	// Auto-migrate the schema
	err = db.AutoMigrate(&model2.User{}, &model2.Category{}, &model2.Tag{}, &model2.Post{}, &model2.SystemSetting{})
	if err != nil {
		log.Printf("Failed to auto-migrate database: %v", err)
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}
	fmt.Println("Database schema migrated successfully.")

	return db, nil
}
