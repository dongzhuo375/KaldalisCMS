package main

import (
	"KaldalisCMS/internal/repository"
	"KaldalisCMS/internal/router"
	"log"
)

func main() {
	// Initialize configuration
	InitConfig()

	// Initialize database
	dsn := GetDatabaseDSN()
	db := repository.InitDB(dsn)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	r := router.SetupRouter()

	log.Println("Server is starting on http://localhost:8080 ...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
