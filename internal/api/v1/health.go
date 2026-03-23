package v1

import (
	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type dbReadinessFunc func(context.Context) error

type HealthAPI struct {
	mode    string
	checkDB dbReadinessFunc
	timeout time.Duration
}

func NewAppHealthAPI(systemSvc *service.SystemService) *HealthAPI {
	return &HealthAPI{
		mode:    "app",
		checkDB: systemSvc.CheckDatabase,
		timeout: 2 * time.Second,
	}
}

func NewSetupHealthAPI() *HealthAPI {
	return &HealthAPI{
		mode:    "setup",
		timeout: 2 * time.Second,
	}
}

func (api *HealthAPI) RegisterRootRoutes(r gin.IRoutes) {
	r.GET("/healthz", api.Healthz)
	r.GET("/readyz", api.Readyz)
}

// Healthz reports process liveness only.
// @Summary Liveness probe
// @Description Returns 200 when process is alive.
// @Tags health
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Router /healthz [get]
func (api *HealthAPI) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, dto.HealthResponse{
		Status: "ok",
		Mode:   api.mode,
		Checks: map[string]dto.HealthCheckResult{},
	})
}

// Readyz reports service readiness for traffic.
// @Summary Readiness probe
// @Description Returns 200 only when dependencies are ready. In setup mode, always returns 503.
// @Tags health
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Failure 503 {object} dto.HealthResponse
// @Router /readyz [get]
func (api *HealthAPI) Readyz(c *gin.Context) {
	if api.mode == "setup" {
		c.JSON(http.StatusServiceUnavailable, dto.HealthResponse{
			Status: "not_ready",
			Mode:   api.mode,
			Checks: map[string]dto.HealthCheckResult{
				"database": {Status: "skip", Detail: "setup mode"},
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), api.timeout)
	defer cancel()

	if api.checkDB == nil || api.checkDB(ctx) != nil {
		c.JSON(http.StatusServiceUnavailable, dto.HealthResponse{
			Status: "not_ready",
			Mode:   api.mode,
			Checks: map[string]dto.HealthCheckResult{
				"database": {Status: "fail", Detail: "database ping failed"},
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.HealthResponse{
		Status: "ok",
		Mode:   api.mode,
		Checks: map[string]dto.HealthCheckResult{
			"database": {Status: "ok"},
		},
	})
}
