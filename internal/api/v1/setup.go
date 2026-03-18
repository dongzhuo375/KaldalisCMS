package v1

import (
	"KaldalisCMS/internal/api/errorx"
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

// Status returns setup mode status.
func (api *SetupAPI) Status(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SystemStatusResponse{Installed: false})
}

// CheckDB validates database connectivity before installation.
// @Summary Check database connection
// @Description Validate DB connection using setup credentials.
// @Tags setup
// @Accept json
// @Produce json
// @Param body body dto.CheckDBRequest true "database connection payload"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /system/check-db [post]
func (api *SetupAPI) CheckDB(c *gin.Context) {
	var req dto.CheckDBRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorx.RespondValidationError(c, "invalid request body", map[string]any{"reason": err.Error()})
		return
	}

	if err := api.svc.ValidateDatabase(req.Host, req.Port, req.User, req.Pass, req.Name); err != nil {
		errorx.RespondErrorByCore(c, err, http.StatusInternalServerError, nil)
		return
	}

	errorx.RespondMessage(c, http.StatusOK, "database connection check passed")
}

// Setup runs first-time installation workflow.
func (api *SetupAPI) Setup(c *gin.Context) {
	var req dto.SystemSetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorx.RespondValidationError(c, "invalid request body", map[string]any{"reason": err.Error()})
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
		errorx.RespondErrorByCore(c, err, http.StatusInternalServerError, nil)
		return
	}

	errorx.RespondMessage(c, http.StatusOK, "installation succeeded, restarting system")
}
