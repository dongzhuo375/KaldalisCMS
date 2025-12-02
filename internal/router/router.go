package router

import (
	v1 "KaldalisCMS/internal/api/v1"
	"KaldalisCMS/internal/core"
	repository "KaldalisCMS/internal/repository/postgres"
	"KaldalisCMS/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {	r := gin.Default()

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

	//Dependency Injection for Post
	var postRepo core.PostRepository = repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postAPI := v1.NewPostAPI(postService)

	 // Dependency Injection for User 
    //var userRepo core.UserRepository = repository.NewUserRepository(db)
     //userService := service.NewUserService(userRepo)
     // userAPI := v1.NewUserAPI(userService) // 下一步就是创建UserAPI并注入userService

	//Register routes
	apiV1 := r.Group("/api/v1")
	postAPI.RegisterRoutes(apiV1)


	return r
}
