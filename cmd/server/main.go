package main

import (
	"KaldalisCMS/internal/infra/auth" // 新增导入
	"KaldalisCMS/internal/infra/repository/postgres"
	"KaldalisCMS/internal/router"
	"log"
)

func main() {
	// Initialize configuration
	InitConfig()

	// Initialize database
	dsn := GetDatabaseDSN()
	db, err := repository.InitDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// --- Casbin 初始化 ---
	// 调用封装好的 Casbin 初始化函数
	enforcer := auth.InitCasbin(db, auth.CasbinConfig{
		ModelPath: "cmd/configs/casbin_model.conf", // Casbin 模型文件路径
	})

	// 规则初始化 (在 main 中)
	// 4. 添加一些初始的权限策略，以便测试
	//    如果策略已存在，则不会重复添加
	// 格式: p, 角色, 路由, 方法
	// 示例：允许 admin 角色访问 /api/v1/posts 并执行 POST 请求
	if has, _ := enforcer.AddPolicy("admin", "/api/v1/posts", "POST"); !has {
		log.Println("策略已存在: admin, /api/v1/posts, POST")
	}
	if has, _ := enforcer.AddPolicy("admin", "/api/v1/posts/:id", "PUT"); !has {
		log.Println("策略已存在: admin, /api/v1/posts/:id, PUT")
	}
	if has, _ := enforcer.AddPolicy("admin", "/api/v1/posts/:id", "DELETE"); !has {
		log.Println("策略已存在: admin, /api/v1/posts/:id, DELETE")
	}
	// 允许所有角色 (包括匿名) 访问 /api/v1/posts 并执行 GET 请求
	if has, _ := enforcer.AddPolicy("anonymous", "/api/v1/posts", "GET"); !has {
		log.Println("策略已存在: anonymous, /api/v1/posts, GET")
	}
	if has, _ := enforcer.AddPolicy("user", "/api/v1/posts", "GET"); !has {
		log.Println("策略已存在: user, /api/v1/posts, GET")
	}
	if has, _ := enforcer.AddPolicy("admin", "/api/v1/posts", "GET"); !has {
		log.Println("策略已存在: admin, /api/v1/posts, GET")
	}

	// --- Casbin 初始化结束 ---

	// 将 enforcer 传递给路由设置函数
	r := router.SetupRouter(db, AppConfig.Auth, enforcer)

	log.Println("Server is starting on http://localhost:8080 ...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
