package v1

import (
	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SetupAPI struct {
	svc *service.SetupService
}

func NewSetupAPI(svc *service.SetupService) *SetupAPI {
	return &SetupAPI{svc: svc}
}

func (api *SetupAPI) RegisterRoutes(r *gin.RouterGroup) {
	system := r.Group("/system")
	{
		system.GET("/status", api.Status)
		system.POST("/setup", api.Setup)
		system.POST("/check-db", api.CheckDB) // 新增：预检接口
	}
}

func (api *SetupAPI) Status(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SystemStatusResponse{Installed: false})
}

func (api *SetupAPI) CheckDB(c *gin.Context) {
	var req struct {
		Host string `json:"host" binding:"required"`
		Port int    `json:"port" binding:"required"`
		User string `json:"user" binding:"required"`
		Pass string `json:"pass"`
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数不完整"})
		return
	}

	if err := api.svc.ValidateDatabase(req.Host, req.Port, req.User, req.Pass, req.Name); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "数据库连接测试通过"})
}

func (api *SetupAPI) Setup(c *gin.Context) {
	var req dto.SystemSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := service.SetupConfig{
		DbHost:             req.DbHost,
		DbPort:             req.DbPort,
		DbUser:             req.DbUser,
		DbPass:             req.DbPass,
		DbName:             req.DbName,
		SiteName:           req.SiteName,
		AdminUser:          req.AdminUsername,
		AdminPass:          req.AdminPassword,
		AdminEmail:         req.AdminEmail,
		AllowAnonymousRead: req.AllowAnonymousRead,
		AdminFullAccess:    req.AdminFullAccess,
	}
	if err := api.svc.Install(cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "安装成功，系统重启中..."})
}
