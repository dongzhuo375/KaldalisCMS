package auth

import (
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// CasbinConfig 定义 Casbin 需要的配置
type CasbinConfig struct {
	ModelPath string // 例如 "configs/rbac_model.conf"
}

// InitCasbin 初始化 Enforcer
// 注意：这里依赖了 *gorm.DB，这是典型的 infra 层依赖
func InitCasbin(db *gorm.DB, cfg CasbinConfig) *casbin.Enforcer {
	// 1. 初始化 Gorm 适配器 (自动建表 casbin_rule)
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		// 这里可以直接 panic，因为鉴权组件起不来，服务也没法跑
		log.Panicf("failed to initialize casbin adapter: %v", err)
	}

	// 2. 加载模型配置文件
	enforcer, err := casbin.NewEnforcer(cfg.ModelPath, adapter)
	if err != nil {
		log.Panicf("failed to create casbin enforcer: %v", err)
	}

	// 3. 开启自动保存策略 (AutoSave)
	// 这样你在代码里 AddPolicy 时，它会自动写库
	enforcer.EnableAutoSave(true)

	// 4. 从数据库加载策略到内存
	if err := enforcer.LoadPolicy(); err != nil {
		log.Panicf("failed to load casbin policy: %v", err)
	}

	return enforcer
}
