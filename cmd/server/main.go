package main

import (
	"KaldalisCMS/internal/router"
	"log"
)

func main() {
	r := router.SetupRouter()

	log.Println("Server is starting on http://localhost:8080 ...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
