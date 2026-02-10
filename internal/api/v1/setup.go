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
	}
}

func (api *SetupAPI) Status(c *gin.Context) {
	// In Setup Mode, system is by definition NOT installed.
	c.JSON(http.StatusOK, dto.SystemStatusResponse{Installed: false})
}

func (api *SetupAPI) Setup(c *gin.Context) {
	var req dto.SystemSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := service.SetupConfig{
		DbHost:     req.DbHost,
		DbPort:     req.DbPort,
		DbUser:     req.DbUser,
		DbPass:     req.DbPass,
		DbName:     req.DbName,
		SiteName:   req.SiteName,
		AdminUser:  req.AdminUsername,
		AdminPass:  req.AdminPassword,
		AdminEmail: req.AdminEmail,
	}

	if err := api.svc.Install(cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setup completed. System is restarting..."})
}
