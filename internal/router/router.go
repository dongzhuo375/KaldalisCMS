package router

import (
	v1 "KaldalisCMS/internal/api/v1"
	"KaldalisCMS/internal/infra/auth"
	repository "KaldalisCMS/internal/infra/repository/postgres"
	"KaldalisCMS/internal/service"

	apimw "KaldalisCMS/internal/api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, authCfg auth.Config) *gin.Engine {
	r := gin.Default()

	// Add a simple CORS middleware
	r.Use(apimw.CORSMiddleware())

	// Dependency Injection for Post

	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postAPI := v1.NewPostAPI(postService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	sessionMgr := auth.NewSessionManager(authCfg)
	userAPI := v1.NewUserAPI(userService, sessionMgr)

	apiV1 := r.Group("/api/v1")
	{
		apiV1.Use(apimw.OptionalAuth(sessionMgr))

		// --- 公开路由 ---
		// 用户登录注册
		userAPI.RegisterRoutes(apiV1)

		// 文章只读接口
		apiV1.GET("/posts", postAPI.GetPosts)
		apiV1.GET("/posts/:id", postAPI.GetPostByID)

		// --- 受保护路由组 ---
		protected := apiV1.Group("/")

		protected.Use(apimw.RequireAuth())
		protected.Use(apimw.CSRFCheck(sessionMgr))

		{
			// 需要认证的 User 操作
			protected.POST("/users/logout", userAPI.Logout)

			// 需要认证的 Post 操作
			protected.POST("/posts", postAPI.CreatePost)
			protected.PUT("/posts/:id", postAPI.UpdatePost)
			protected.DELETE("/posts/:id", postAPI.DeletePost)
		}
	}

	return r
}
