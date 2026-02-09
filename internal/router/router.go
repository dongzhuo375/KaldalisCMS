package router

import (
	apimw "KaldalisCMS/internal/api/middleware"
	v1 "KaldalisCMS/internal/api/v1"
	"KaldalisCMS/internal/infra/auth"
	repository "KaldalisCMS/internal/infra/repository/postgres"
	"KaldalisCMS/internal/service"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewAppRouter initializes the router for the fully functional application
func NewAppRouter(db *gorm.DB, authCfg auth.Config, enforcer *casbin.Enforcer) *gin.Engine {
	r := gin.Default()

	// Add a simple CORS middleware
	r.Use(apimw.CORSMiddleware())

	// Dependency Injection
	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postAPI := v1.NewPostAPI(postService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	sessionMgr := auth.NewSessionManager(authCfg)
	userAPI := v1.NewUserAPI(userService, sessionMgr)

	systemRepo := repository.NewSystemRepository(db)
	systemService := service.NewSystemService(db, systemRepo, userService)
	systemAPI := v1.NewSystemAPI(systemService)

	apiV1 := r.Group("/api/v1")
	{
		apiV1.Use(apimw.OptionalAuth(sessionMgr))

		userAPI.RegisterRoutes(apiV1)
		systemAPI.RegisterRoutes(apiV1)

		apiV1.GET("/posts", postAPI.GetPosts)
		apiV1.GET("/posts/:id", postAPI.GetPostByID)

		protected := apiV1.Group("/")
		protected.Use(apimw.RequireAuth())
		protected.Use(apimw.Authorize(enforcer))
		protected.Use(apimw.CSRFCheck(sessionMgr))

		{
			protected.POST("/users/logout", userAPI.Logout)
			protected.POST("/posts", postAPI.CreatePost)
			protected.PUT("/posts/:id", postAPI.UpdatePost)
			protected.DELETE("/posts/:id", postAPI.DeletePost)
		}
	}

	return r
}

// NewSetupRouter initializes the router for the setup mode, hiding service instantiation from main
func NewSetupRouter(save func(string, int, string, string, string) error, reload func() error) *gin.Engine {
	r := gin.Default()
	r.Use(apimw.CORSMiddleware())

	// Service instantiation is now hidden inside router
	setupSvc := service.NewSetupService(save, reload)
	setupAPI := v1.NewSetupAPI(setupSvc)
	
	apiV1 := r.Group("/api/v1")
	{
		setupAPI.RegisterRoutes(apiV1)
	}

	return r
}
