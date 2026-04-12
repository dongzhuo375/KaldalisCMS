package router

import (
	apimw "KaldalisCMS/internal/api/middleware"
	v1 "KaldalisCMS/internal/api/v1"
	"KaldalisCMS/internal/infra/auth"
	repository "KaldalisCMS/internal/infra/repository/postgres"
	"KaldalisCMS/internal/service"
	"KaldalisCMS/internal/utils"

	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

// ensurePostWorkflowPolicies backfills the complete set of post-management
// policies required by the current application build.  It covers:
//   - role inheritance
//   - route policies for the HTTP authorization middleware
//   - capability policies for the service-layer post authorizer
//
// AddPolicy / AddGroupingPolicy are idempotent — duplicates are silently ignored.
func ensurePostWorkflowPolicies(enforcer *casbin.Enforcer) {
	if enforcer == nil {
		return
	}

	// 1. Role inheritance (fix stale reverse inheritance if present)
	_, _ = enforcer.RemoveGroupingPolicy("admin", "super_admin")
	_, _ = enforcer.AddGroupingPolicy("super_admin", "admin")
	_, _ = enforcer.AddGroupingPolicy("admin", "user")

	// 2. Route policies — must mirror what setup_service seeds
	rules := [][]string{
		// logout (all authenticated roles)
		{"admin", "/api/v1/users/logout", "POST"},
		{"user", "/api/v1/users/logout", "POST"},
		{"super_admin", "/api/v1/users/logout", "POST"},

		// admin route policies
		{"admin", "/api/v1/admin/posts", "GET"},
		{"admin", "/api/v1/admin/posts", "POST"},
		{"admin", "/api/v1/admin/posts/:id", "GET"},
		{"admin", "/api/v1/admin/posts/:id", "PUT"},
		{"admin", "/api/v1/admin/posts/:id/publish", "POST"},
		{"admin", "/api/v1/admin/posts/:id/draft", "POST"},

		// admin capability policies
		{"admin", "post", "list:any"},
		{"admin", "post", "read:any"},
		{"admin", "post", "update:any"},
		{"admin", "post", "publish"},
		{"admin", "post", "unpublish"},
		{"admin", "post", "delete"},

		// user route policies
		{"user", "/api/v1/posts", "GET"},
		{"user", "/api/v1/posts/:id", "GET"},
		{"user", "/api/v1/admin/posts", "GET"},
		{"user", "/api/v1/admin/posts", "POST"},
		{"user", "/api/v1/admin/posts/:id", "GET"},
		{"user", "/api/v1/admin/posts/:id", "PUT"},
		{"user", "/api/v1/media", "GET"},

		// user capability policies
		{"user", "post:draft", "create"},
		{"user", "post:draft", "list:own"},
		{"user", "post:draft", "read:own"},
		{"user", "post:draft", "update:own"},
	}

	for _, rule := range rules {
		_, _ = enforcer.AddPolicy(rule[0], rule[1], rule[2])
	}

	_ = enforcer.SavePolicy()
}

// NewAppRouter initializes the router for the fully functional application.
func NewAppRouter(db *gorm.DB, authCfg auth.Config, enforcer *casbin.Enforcer, swaggerOpts SwaggerOptions) *gin.Engine {
	r := gin.Default()
	r.Use(apimw.CORSMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	registerSwaggerRoutes(r, swaggerOpts)

	uploadDir := os.Getenv("MEDIA_UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = filepath.FromSlash("./data/uploads")
	}
	maxUploadMB := os.Getenv("MEDIA_MAX_UPLOAD_SIZE_MB")
	publicBaseURL := os.Getenv("MEDIA_PUBLIC_BASE_URL")
	maxFilenameBytes := os.Getenv("MEDIA_MAX_FILENAME_BYTES")

	r.Static("/media/a", filepath.Join(uploadDir, "a"))

	mediaRepo := repository.NewMediaRepository(db)
	mediaCfg := service.MediaConfig{UploadDir: uploadDir}
	if v := utils.ParseInt64(maxUploadMB); v > 0 {
		mediaCfg.MaxUploadSizeMB = v
	} else {
		mediaCfg.MaxUploadSizeMB = 50
	}
	mediaCfg.PublicBaseURL = publicBaseURL
	if v := utils.ParseInt(maxFilenameBytes); v > 0 {
		mediaCfg.MaxFilenameBytes = v
	} else {
		mediaCfg.MaxFilenameBytes = 180
	}
	mediaSvc := service.NewMediaService(mediaRepo, mediaCfg)
	mediaAPI := v1.NewMediaAPI(mediaSvc, mediaRepo)

	postRepo := repository.NewPostRepository(db)
	postAuthorizer := auth.NewCasbinPostAuthorizer(enforcer)
	postService := service.NewPostServiceWithMedia(postRepo, mediaSvc, postAuthorizer)
	publicPostAPI := v1.NewPublicPostAPI(postService)
	adminPostAPI := v1.NewAdminPostAPI(postService)
	ensurePostWorkflowPolicies(enforcer)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	sessionMgr := auth.NewSessionManager(authCfg)
	userAPI := v1.NewUserAPI(userService, sessionMgr)

	systemRepo := repository.NewSystemRepository(db)
	systemService := service.NewSystemService(db, systemRepo, userService)
	systemAPI := v1.NewSystemAPI(systemService)
	healthAPI := v1.NewAppHealthAPI(systemService)
	healthAPI.RegisterRootRoutes(r) // Correct: root registration

	go func() {
		utils.RunTicker(1*time.Hour, func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()
			if err := mediaSvc.CleanupStaleMedia(ctx); err != nil {
				println("Error cleaning up media:", err.Error())
			}
		})
	}()

	apiV1 := r.Group("/api/v1")
	apiV1.Use(apimw.OptionalAuth(sessionMgr))
	{
		userAPI.RegisterRoutes(apiV1)
		systemAPI.RegisterRoutes(apiV1)

		// Public post routes go through Casbin so the AllowAnonymousRead
		// setting actually controls anonymous access.  OptionalAuth is
		// already applied on the parent group; Authorize falls back to
		// "anonymous" when no role is present in the context.
		public := apiV1.Group("/")
		public.Use(apimw.Authorize(enforcer))
		{
			public.GET("/posts", publicPostAPI.GetPosts)
			public.GET("/posts/:id", publicPostAPI.GetPostByID)
		}

		protected := apiV1.Group("/")
		protected.Use(apimw.RequireAuth())
		protected.Use(apimw.Authorize(enforcer))
		protected.Use(apimw.CSRFCheck(sessionMgr))
		{
			protected.POST("/users/logout", userAPI.Logout)

			adminPosts := protected.Group("/admin")
			adminPosts.GET("/posts", adminPostAPI.GetPosts)
			adminPosts.GET("/posts/:id", adminPostAPI.GetPostByID)
			adminPosts.POST("/posts", adminPostAPI.CreatePost)
			adminPosts.PUT("/posts/:id", adminPostAPI.UpdatePost)
			adminPosts.POST("/posts/:id/publish", adminPostAPI.PublishPost)
			adminPosts.POST("/posts/:id/draft", adminPostAPI.DraftPost)
			adminPosts.DELETE("/posts/:id", adminPostAPI.DeletePost)

			mediaAPI.RegisterRoutes(protected)
		}
	}

	return r
}

func NewSetupRouter(save func(string, int, string, string, string) error, reload func() error, swaggerOpts SwaggerOptions) *gin.Engine {
	r := gin.Default()
	r.Use(apimw.CORSMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	registerSwaggerRoutes(r, swaggerOpts)

	setupSvc := service.NewSetupService(save, reload)
	setupAPI := v1.NewSetupAPI(setupSvc)
	healthAPI := v1.NewSetupHealthAPI()
	healthAPI.RegisterRootRoutes(r) // Correct: root registration

	apiV1 := r.Group("/api/v1")
	setupAPI.RegisterRoutes(apiV1)

	return r
}
