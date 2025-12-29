package v1

import (
	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserAPI struct {
	service core.UserService
	sm      core.SessionManager
}

func NewUserAPI(service core.UserService, sessionMgr core.SessionManager) *UserAPI {
	return &UserAPI{
		service: service,
		sm:      sessionMgr,
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
	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser := entity.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     "user", // Assign default role
	}

	if err := api.service.CreateUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (a *UserAPI) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := a.service.Login(ctx, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err := a.sm.EstablishSession(c.Writer, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登录状态创建失败"})
		return
	}

	// 动态计算过期时间
	expiresAt := time.Now().Add(a.sm.GetTTL())

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
		"expires_at": expiresAt.Format(time.RFC3339), // 可改为从 manager 读取
	})
}

func (a *UserAPI) Logout(c *gin.Context) {
	// Logout 通过 service 层触发副作用
	//a.service.Logout() 暂时无逻辑
	a.sm.DestroySession(c.Writer)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
