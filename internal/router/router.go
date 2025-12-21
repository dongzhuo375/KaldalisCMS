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
	r.Use(apimw.CORSMiddleware())

	// API V1 Group
	apiV1 := r.Group("/api/v1")

		// Dependency Injection for Post

		var postRepo core.PostRepository = repository.NewPostRepository(db)

		postService := service.NewPostService(postRepo)

		postAPI := v1.NewPostAPI(postService)

	

		// Public post routes (read-only) are registered on the main v1 group

		apiV1.GET("/posts", postAPI.GetPosts)

		apiV1.GET("/posts/:id", postAPI.GetPostByID)

	

		// Dependency Injection for User

		var userRepo core.UserRepository = repository.NewUserRepository(db)

	

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

		// keep existing service construction (do not change business logic)

		userService := service.NewUserService(userRepo, authMgr)

	

		// Ensure OptionalAuthWithManager is actually used: apply it to apiV1 so handlers

		// can observe optional login state (useful for SSR/public endpoints).

		apiV1.Use(apimw.OptionalAuthWithManager(authMgr))

	

		// Construct UserAPI with existing userService and the singleton auth manager.

		// This keeps your original user service logic while enabling infra-level cookie/jwt/csrf handling.

		userAPI := v1.NewUserAPI(userService)

		userAPI.RegisterRoutes(apiV1) // Register User routes (register/login only)

	

		// Protected group: require authentication + CSRF for mutating operations.

		protected := apiV1.Group("/")

		protected.Use(apimw.RequireAuthWithManager(authMgr))

		protected.Use(apimw.CSRFMiddlewareWithManager(authMgr))

	

		// Move logout into protected group so only authenticated requests with valid CSRF can logout.

		protected.POST("/users/logout", userAPI.Logout)

	

		// Add protected post routes for writing/mutating data

		protected.POST("/posts", postAPI.CreatePost)

		protected.PUT("/posts/:id", postAPI.UpdatePost)

		protected.DELETE("/posts/:id", postAPI.DeletePost)

	_ = apimw.RequireAuthWithManager // keep linter happy if middlewares unused elsewhere

	return r
}
