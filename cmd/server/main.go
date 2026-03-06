package main

import (
	"KaldalisCMS/internal/infra/auth"
	"KaldalisCMS/internal/infra/model"
	repository "KaldalisCMS/internal/infra/repository/postgres"
	"KaldalisCMS/internal/router"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

// RouterManager acts as a dynamic proxy for the active http.Handler
type RouterManager struct {
	mu      sync.RWMutex
	current http.Handler
}

func (rm *RouterManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rm.mu.RLock()
	handler := rm.current
	rm.mu.RUnlock()
	if handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Service initializing...", http.StatusServiceUnavailable)
	}
}

func (rm *RouterManager) Switch(h http.Handler) {
	rm.mu.Lock()
	rm.current = h
	rm.mu.Unlock()
}

var routerManager = &RouterManager{}

func main() {
	// Initialize configuration
	InitConfig()

	// Try to bootstrap the full application
	if err := BootstrapApp(); err != nil {
		log.Printf("系统启动检查未通过: %v. 切换到 [安装模式] (SETUP MODE).", err)
		SwitchToSetupMode()
	}

	log.Println("服务器正在启动，监听端口: http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", routerManager); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// BootstrapApp 不仅测试连接，还会检查数据库是否完整且标记为已安装
func BootstrapApp() error {
	dsn := GetDatabaseDSN()

	db, err := repository.InitDB(dsn)
	if err != nil {
		return err
	}

	// --- 核心：检查安装状态 ---
	var setting model.SystemSetting
	// 尝试读取 ID 为 1 的系统设置
	if err := db.First(&setting, 1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("系统尚未初始化 (数据库记录缺失)")
		}
		return err
	}
	if !setting.Installed {
		return errors.New("系统已连接但标记为 [未安装]")
	}

	// --- Casbin 初始化 ---
	enforcer := auth.InitCasbin(db, auth.CasbinConfig{
		ModelPath: "cmd/configs/casbin_model.conf",
	})

	// 初始化策略 (如果不存在)
	setupPolicies(enforcer)

	// --- 启动应用路由 ---
	r := router.NewAppRouter(db, AppConfig.Auth, enforcer)

	routerManager.Switch(r)
	log.Printf("系统正常运行中 [业务模式] (APP MODE) - 站点名称: %s", setting.SiteName)
	return nil
}

func SwitchToSetupMode() {
	r := router.NewSetupRouter(
		SaveDatabaseConfig,
		func() error {
			log.Println("配置已保存，正在尝试热重启进入业务模式...")
			if err := BootstrapApp(); err != nil {
				log.Printf("热重启失败 (安装可能未完全成功): %v", err)
				return err
			}
			log.Println("热重启成功，系统已进入业务模式!")
			return nil
		},
	)

	routerManager.Switch(r)
	log.Println("!!! 系统当前处于 [安装模式] (SETUP MODE) !!!")
}

func setupPolicies(enforcer *casbin.Enforcer) {
	// 管理员权限
	enforcer.AddPolicy("admin", "/api/v1/posts", "POST")
	enforcer.AddPolicy("admin", "/api/v1/posts/:id", "PUT")
	enforcer.AddPolicy("admin", "/api/v1/posts/:id", "DELETE")
	enforcer.AddPolicy("admin", "/api/v1/media", "POST")
	enforcer.AddPolicy("admin", "/api/v1/media", "GET")
	enforcer.AddPolicy("admin", "/api/v1/media/:id", "DELETE")
	enforcer.AddPolicy("admin", "/api/v1/posts/:id/media", "GET")

	// 普通用户与访客权限
	enforcer.AddPolicy("anonymous", "/api/v1/posts", "GET")
	enforcer.AddPolicy("user", "/api/v1/posts", "GET")
	enforcer.AddPolicy("user", "/api/v1/media", "POST")
	enforcer.AddPolicy("user", "/api/v1/media", "GET")
	enforcer.AddPolicy("user", "/api/v1/posts/:id/media", "GET")
}
