package router

import (
	v1 "KaldalisCMS/internal/api/v1"
	"KaldalisCMS/internal/core"
	repository "KaldalisCMS/internal/infra/repository/postgres" // Corrected import alias
	"KaldalisCMS/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, jwtSecret string, jwtExpHours int) *gin.Engine {
	r := gin.Default()

	// Add a simple CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API V1 Group
	apiV1 := r.Group("/api/v1")

	// Dependency Injection for Post
	var postRepo core.PostRepository = repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postAPI := v1.NewPostAPI(postService)
	postAPI.RegisterRoutes(apiV1) // Register Post routes

	// Dependency Injection for User
	var userRepo core.UserRepository = repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtSecret, jwtExpHours)
	userAPI := v1.NewUserAPI(userService) // NewUserAPI now takes concrete *service.UserService
	userAPI.RegisterRoutes(apiV1)         // Register User routes

	return r
}

