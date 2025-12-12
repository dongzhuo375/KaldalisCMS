package router

import (
	v1 "KaldalisCMS/internal/api/v1"
	"KaldalisCMS/internal/core"
	repository "KaldalisCMS/internal/infra/repository/postgres"
	"KaldalisCMS/internal/service"

	apimw "KaldalisCMS/internal/api/middleware"
	authinfra "KaldalisCMS/internal/infra/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"net/http"
	"time"
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
	// keep existing service construction (do not change business logic)
	userService := service.NewUserService(userRepo)

	// Create a single auth manager instance (singleton) and pass it into UserAPI and middlewares.
	// Cookie names, path, domain, secure and sameSite can be adjusted or moved to config as needed.
	authMgr := authinfra.NewManager(
		jwtSecret,
		time.Duration(jwtExpHours)*time.Hour,
		"kaldalis_auth",
		"kaldalis_csrf",
		"/",
		"",   // domain
		true, // secure by default; in local dev you can set to false via config
		http.SameSiteLaxMode,
	)

	// Ensure OptionalAuthWithManager is actually used: apply it to apiV1 so handlers
	// can observe optional login state (useful for SSR/public endpoints).
	apiV1.Use(apimw.OptionalAuthWithManager(authMgr))

	// Construct UserAPI with existing userService and the singleton auth manager.
	// This keeps your original user service logic while enabling infra-level cookie/jwt/csrf handling.
	userAPI := v1.NewUserAPI(userService, authMgr)
	userAPI.RegisterRoutes(apiV1) // Register User routes (register/login only)

	// Protected group: require authentication + CSRF for mutating operations.
	protected := apiV1.Group("/")
	protected.Use(apimw.RequireAuthWithManager(authMgr))
	protected.Use(apimw.CSRFMiddlewareWithManager(authMgr))

	// Move logout into protected group so only authenticated requests with valid CSRF can logout.
	protected.POST("/users/logout", userAPI.Logout)

	// Example: protect post write endpoints (if postAPI exposes handlers). If your postAPI's RegisterRoutes
	// currently registered create/update/delete publicly, you should move those registrations into protected.
	// If postAPI exposes handler functions, register them here. Example (uncomment and adapt if available):
	// protected.POST("/posts", postAPI.Create)
	// protected.PUT("/posts/:id", postAPI.Update)
	// protected.DELETE("/posts/:id", postAPI.Delete)

	_ = apimw.RequireAuthWithManager // keep linter happy if middlewares unused elsewhere

	return r
}
