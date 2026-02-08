package v1

import (
	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SystemAPI struct {
	svc *service.SystemService
}

func NewSystemAPI(svc *service.SystemService) *SystemAPI {
	return &SystemAPI{svc: svc}
}

func (api *SystemAPI) RegisterRoutes(r *gin.RouterGroup) {
	system := r.Group("/system")
	{
		system.GET("/status", api.Status)
		system.POST("/setup", api.Setup)
	}
}

func (api *SystemAPI) Status(c *gin.Context) {
	st, err := api.svc.Status(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.SystemStatusResponse{Installed: st.Installed, SiteName: st.SiteName})
}

func (api *SystemAPI) Setup(c *gin.Context) {
	var req dto.SystemSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := api.svc.SetupOnce(c.Request.Context(), service.SetupParams{
		SiteName:      req.SiteName,
		AdminUsername: req.AdminUsername,
		AdminEmail:    req.AdminEmail,
		AdminPassword: req.AdminPassword,
	})
	if err != nil {
		if errors.Is(err, service.ErrAlreadyInstalled) {
			c.JSON(http.StatusConflict, gin.H{"error": "already installed"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "setup completed"})
}
