package v1

import (
	"KaldalisCMS/internal/core/entity"
	"KaldalisCMS/internal/infra/auth"
	"KaldalisCMS/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserAPI struct {
	service *service.UserService
	authMgr *auth.Manager
}

func NewUserAPI(service *service.UserService, authMgr *auth.Manager) *UserAPI {
	return &UserAPI{
		service: service,
		authMgr: authMgr,
	}
}

// RegisterRoutes registers the user-related routes to the Gin router.
func (api *UserAPI) RegisterRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/register", api.Register)
		userRoutes.POST("/login", api.Login)
	}
}

// Register handles new user registration.
func (api *UserAPI) Register(c *gin.Context) {
	var newUser entity.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := api.service.CreateUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (a *UserAPI) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.service.VerifyUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// secureFlag 可基于配置或请求 TLS 决定
	secureFlag := c.Request.TLS != nil
	if err := a.authMgr.Login(c.Writer, user.ID, secureFlag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set auth cookies"})
		return
	}

	// 返回用户信息（不返回 token）
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
		"expires_at": time.Now().Add(a.authMgr.TTL).Format(time.RFC3339),
	})
}

func (a *UserAPI) Logout(c *gin.Context) {
	a.authMgr.Logout(c.Writer)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
