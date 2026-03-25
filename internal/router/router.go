package router

import (
	apimw "KaldalisCMS/internal/api/middleware"
	v1 "KaldalisCMS/internal/api/v1"
	"KaldalisCMS/internal/infra/auth"
	repository "KaldalisCMS/internal/infra/repository/postgres"
	"KaldalisCMS/internal/service"
	"KaldalisCMS/internal/utils"

	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

// ensurePostWorkflowPolicies backfills the post-management policies required by
// the current application build. It intentionally seeds two complementary policy sets:
//   - route policies used by the HTTP authorization middleware
//   - capability policies used by the service-layer post authorizer

// The ensurePostWorkflowPolicies function silently seeds policies during application startup without checking if they already exist.
// The AddPolicies method on line 68 may fail if policies already exist (depending on Casbin configuration),
// but the error is only logged as a warning. Consider checking for existing policies first, or documenting that this function is idempotent and safe to call on every startup.
func ensurePostWorkflowPolicies(enforcer *casbin.Enforcer) {
	if enforcer == nil {
		return
	}

	rules := [][]string{
		{"user", "/api/v1/admin/posts", "GET"},
		{"user", "/api/v1/admin/posts", "POST"},
		{"user", "/api/v1/admin/posts/:id", "GET"},
		{"user", "/api/v1/admin/posts/:id", "PUT"},
		{"user", "post:draft", "create"},
		{"user", "post:draft", "list:own"},
		{"user", "post:draft", "read:own"},
		{"user", "post:draft", "update:own"},
		{"admin", "/api/v1/admin/posts", "GET"},
		{"admin", "/api/v1/admin/posts", "POST"},
		{"admin", "/api/v1/admin/posts/:id", "GET"},
		{"admin", "/api/v1/admin/posts/:id", "PUT"},
		{"admin", "/api/v1/admin/posts/:id", "DELETE"},
		{"admin", "/api/v1/admin/posts/:id/publish", "POST"},
		{"admin", "/api/v1/admin/posts/:id/draft", "POST"},
		{"admin", "post", "list:any"},
		{"admin", "post", "read:any"},
		{"admin", "post", "update:any"},
		{"admin", "post", "publish"},
		{"admin", "post", "unpublish"},
		{"admin", "post", "delete"},
		{"super_admin", "/api/v1/admin/posts", "GET"},
		{"super_admin", "/api/v1/admin/posts", "POST"},
		{"super_admin", "/api/v1/admin/posts/:id", "GET"},
		{"super_admin", "/api/v1/admin/posts/:id", "PUT"},
		{"super_admin", "/api/v1/admin/posts/:id", "DELETE"},
		{"super_admin", "/api/v1/admin/posts/:id/publish", "POST"},
		{"super_admin", "/api/v1/admin/posts/:id/draft", "POST"},
		{"super_admin", "post", "list:any"},
		{"super_admin", "post", "read:any"},
		{"super_admin", "post", "update:any"},
		{"super_admin", "post", "publish"},
		{"super_admin", "post", "unpublish"},
		{"super_admin", "post", "delete"},
	}

	if _, err := enforcer.AddPolicies(rules); err != nil {
		log.Printf("[WARN] failed to ensure post workflow policies: %v", err)
	}
}

// NewAppRouter initializes the router for the fully functional application.
// Public content delivery and admin management are registered under different paths so
// callers can reason about visibility and authorization from the URL contract alone.
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
	healthAPI.RegisterRootRoutes(r)

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

		// Public post endpoints are permanently read-only and only expose published content.
		apiV1.GET("/posts", publicPostAPI.GetPosts)
		apiV1.GET("/posts/:id", publicPostAPI.GetPostByID)

		protected := apiV1.Group("/")
		protected.Use(apimw.RequireAuth())
		protected.Use(apimw.Authorize(enforcer))
		protected.Use(apimw.CSRFCheck(sessionMgr))
		{
			protected.POST("/users/logout", userAPI.Logout)

			// Management post endpoints stay under /api/v1/admin/posts to keep draft visibility
			// and write operations separate from the public content contract. Casbin still guards
			// route entry here, while the service layer asks a Casbin-backed authorizer for the
			// finer-grained post capabilities needed to resolve own-draft vs any-post access.
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

// NewSetupRouter initializes the router for the setup mode, hiding service instantiation from main.
func NewSetupRouter(save func(string, int, string, string, string) error, reload func() error, swaggerOpts SwaggerOptions) *gin.Engine {
	r := gin.Default()
	r.Use(apimw.CORSMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	registerSwaggerRoutes(r, swaggerOpts)

	setupSvc := service.NewSetupService(save, reload)
	setupAPI := v1.NewSetupAPI(setupSvc)
	healthAPI := v1.NewSetupHealthAPI()
	healthAPI.RegisterRootRoutes(r)

	apiV1 := r.Group("/api/v1")
	setupAPI.RegisterRoutes(apiV1)

	return r
}
