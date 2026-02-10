package router

import (
	apimw "KaldalisCMS/internal/api/middleware"
	v1 "KaldalisCMS/internal/api/v1"
	"KaldalisCMS/internal/infra/auth"
	repository "KaldalisCMS/internal/infra/repository/postgres"
	"KaldalisCMS/internal/service"
	"KaldalisCMS/internal/utils"

	"os"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, authCfg auth.Config, enforcer *casbin.Enforcer) *gin.Engine {
	r := gin.Default()

	// Add a simple CORS middleware
	r.Use(apimw.CORSMiddleware())

	// --- Static media (public) ---
	uploadDir := os.Getenv("MEDIA_UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = filepath.FromSlash("./data/uploads")
	}
	maxUploadMB := os.Getenv("MEDIA_MAX_UPLOAD_SIZE_MB")
	publicBaseURL := os.Getenv("MEDIA_PUBLIC_BASE_URL")
	maxFilenameBytes := os.Getenv("MEDIA_MAX_FILENAME_BYTES")

	// /media maps to uploadDir; assets are stored under uploadDir/a/{id}/{stored_name}
	r.Static("/media", uploadDir)

	// --- Media ---
	mediaRepo := repository.NewMediaRepository(db)
	mediaCfg := service.MediaConfig{
		UploadDir: uploadDir,
	}
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

	// Dependency Injection for Post

	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostServiceWithMedia(postRepo, mediaSvc)
	postAPI := v1.NewPostAPI(postService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	sessionMgr := auth.NewSessionManager(authCfg)
	userAPI := v1.NewUserAPI(userService, sessionMgr)

	// --- System init/setup ---
	systemRepo := repository.NewSystemRepository(db)
	systemService := service.NewSystemService(db, systemRepo, userService)
	systemAPI := v1.NewSystemAPI(systemService)

	apiV1 := r.Group("/api/v1")
	{
		apiV1.Use(apimw.OptionalAuth(sessionMgr))

		// --- 公开路由 ---
		// 用户登录注册
		userAPI.RegisterRoutes(apiV1)

		// 系统状态与初始化
		systemAPI.RegisterRoutes(apiV1)

		// 文章只读接口
		apiV1.GET("/posts", postAPI.GetPosts)
		apiV1.GET("/posts/:id", postAPI.GetPostByID)

		// --- 受保护路由组 ---
		protected := apiV1.Group("/")

		// 先检查是否登录，再检查 CSRF，最后检查权限
		protected.Use(apimw.RequireAuth())
		protected.Use(apimw.Authorize(enforcer))
		protected.Use(apimw.CSRFCheck(sessionMgr))

		{
			// 需要认证的 User 操作
			protected.POST("/users/logout", userAPI.Logout)

			// 需要认证的 Post 操作
			protected.POST("/posts", postAPI.CreatePost)
			protected.PUT("/posts/:id", postAPI.UpdatePost)
			protected.DELETE("/posts/:id", postAPI.DeletePost)

			// 需要认证的 Media 操作
			mediaAPI.RegisterRoutes(protected)
		}
	}

	return r
}
