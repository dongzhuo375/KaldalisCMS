package main

import (
	"KaldalisCMS/internal/infra/auth"
	"KaldalisCMS/internal/infra/repository/postgres"
	"KaldalisCMS/internal/router"
	"log"
	"net/http"
	"sync"
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
		log.Printf("Initialization failed: %v. Switching to SETUP MODE.", err)
		SwitchToSetupMode()
	}

	log.Println("Server is starting on http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", routerManager); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

// BootstrapApp tries to connect to DB and setup the full application router
func BootstrapApp() error {
	dsn := GetDatabaseDSN()
	
	db, err := repository.InitDB(dsn)
	if err != nil {
		return err
	}

	enforcer := auth.InitCasbin(db, auth.CasbinConfig{
		ModelPath: "cmd/configs/casbin_model.conf",
	})

	// Setup Default Policies
	enforcer.AddPolicy("admin", "/api/v1/posts", "POST")
	enforcer.AddPolicy("admin", "/api/v1/posts/:id", "PUT")
	enforcer.AddPolicy("admin", "/api/v1/posts/:id", "DELETE")
	enforcer.AddPolicy("anonymous", "/api/v1/posts", "GET")
	enforcer.AddPolicy("user", "/api/v1/posts", "GET")
	enforcer.AddPolicy("admin", "/api/v1/posts", "GET")

	// Create App Router
	r := router.NewAppRouter(db, AppConfig.Auth, enforcer)

	// Switch global handler
	routerManager.Switch(r)
	log.Println("System running in APP MODE")
	return nil
}

// SwitchToSetupMode initializes the limited setup router
func SwitchToSetupMode() {
	// Now main only communicates with router, passing callbacks
	// No need to import internal/service here
	r := router.NewSetupRouter(
		SaveDatabaseConfig, // The save callback
		func() error {      // The reload callback
			log.Println("Configuration saved. Attempting hot reload...")
			if err := BootstrapApp(); err != nil {
				log.Printf("Hot reload failed: %v", err)
				return err
			}
			log.Println("Hot reload successful!")
			return nil
		},
	)

	routerManager.Switch(r)
	log.Println("System running in SETUP MODE")
}
