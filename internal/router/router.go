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

// ensurePostWorkflowPolicies backfills the post-management policies required by
// the current application build. It intentionally seeds two complementary policy sets:
//   - route policies used by the HTTP authorization middleware
//   - capability policies used by the service-layer post authorizer
func ensurePostWorkflowPolicies(enforcer *casbin.Enforcer) {
	if enforcer == nil {
		return
	}

	// 1. 建立角色继承关系 (super_admin 继承 admin)
	// 确保没有反向继承
	_, _ = enforcer.RemoveGroupingPolicy("admin", "super_admin")
	_, _ = enforcer.AddGroupingPolicy("super_admin", "admin")

	// 2. 检查并添加缺失的基础策略
	rules := [][]string{
		{"admin", "/api/v1/users/logout", "POST"},
		{"user", "/api/v1/users/logout", "POST"},
		{"super_admin", "/api/v1/users/logout", "POST"},
		
		// Post相关
		{"admin", "post", "delete"},
		{"admin", "post", "publish"},
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

		apiV1.GET("/posts", publicPostAPI.GetPosts)
		apiV1.GET("/posts/:id", publicPostAPI.GetPostByID)

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
