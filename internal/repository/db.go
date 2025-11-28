package repository

import (
	"KaldalisCMS/internal/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB(dsn string) *gorm.DB {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection established successfully.")

	// Auto-migrate the schema
	err = DB.AutoMigrate(&model.User{}, &model.Category{}, &model.Tag{}, &model.Post{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	fmt.Println("Database schema migrated successfully.")

	return DB
}