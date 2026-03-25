package v1

import (
	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/service"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type dbReadinessFunc func(context.Context) error

type HealthAPI struct {
	mode             string
	checkDB          dbReadinessFunc
	timeout          time.Duration
	readySuccessTTL  time.Duration
	readyFailureTTL  time.Duration
	cacheMu          sync.RWMutex
	readyCacheExpire time.Time
	readyCacheCode   int
	readyCacheBody   dto.HealthResponse
	readyCacheValid  bool
}

var probeRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "kaldalis",
		Subsystem: "probe",
		Name:      "requests_total",
		Help:      "Total number of probe requests by probe, mode, result and cache status.",
	},
	[]string{"probe", "mode", "result", "cache"},
)

var probeReadyState = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "kaldalis",
		Subsystem: "probe",
		Name:      "ready_state",
		Help:      "Current readiness state by mode. 1 means ready, 0 means not ready.",
	},
	[]string{"mode"},
)

func init() {
	prometheus.MustRegister(probeRequestsTotal, probeReadyState)
}

func NewAppHealthAPI(systemSvc *service.SystemService) *HealthAPI {
	return &HealthAPI{
		mode:            "app",
		checkDB:         systemSvc.CheckDatabase,
		timeout:         2 * time.Second,
		readySuccessTTL: 400 * time.Millisecond,
		readyFailureTTL: 250 * time.Millisecond,
	}
}

func NewSetupHealthAPI() *HealthAPI {
	return &HealthAPI{
		mode:            "setup",
		timeout:         2 * time.Second,
		readySuccessTTL: 400 * time.Millisecond,
		readyFailureTTL: 250 * time.Millisecond,
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
	probeRequestsTotal.WithLabelValues("healthz", api.mode, "ok", "none").Inc()
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
	if code, body, ok := api.cachedReadyz(); ok {
		probeRequestsTotal.WithLabelValues("readyz", api.mode, body.Status, "hit").Inc()
		probeReadyState.WithLabelValues(api.mode).Set(readyStateFromCode(code))
		c.JSON(code, body)
		return
	}

	code, body := api.evaluateReadyz(c.Request.Context())
	api.cacheReadyz(code, body)

	probeRequestsTotal.WithLabelValues("readyz", api.mode, body.Status, "miss").Inc()
	probeReadyState.WithLabelValues(api.mode).Set(readyStateFromCode(code))
	c.JSON(code, body)
}

func (api *HealthAPI) cachedReadyz() (int, dto.HealthResponse, bool) {
	now := time.Now()

	api.cacheMu.RLock()
	code := api.readyCacheCode
	body := api.readyCacheBody
	expiresAt := api.readyCacheExpire
	valid := api.readyCacheValid
	api.cacheMu.RUnlock()

	if !valid || now.After(expiresAt) {
		return 0, dto.HealthResponse{}, false
	}
	return code, body, true
}

func (api *HealthAPI) cacheReadyz(code int, body dto.HealthResponse) {
	ttl := api.readySuccessTTL
	if code != http.StatusOK {
		ttl = api.readyFailureTTL
	}
	if ttl <= 0 {
		return
	}

	api.cacheMu.Lock()
	api.readyCacheCode = code
	api.readyCacheBody = body
	api.readyCacheExpire = time.Now().Add(ttl)
	api.readyCacheValid = true
	api.cacheMu.Unlock()
}

func (api *HealthAPI) evaluateReadyz(ctx context.Context) (int, dto.HealthResponse) {
	if api.mode == "setup" {
		return http.StatusServiceUnavailable, dto.HealthResponse{
			Status: "not_ready",
			Mode:   api.mode,
			Checks: map[string]dto.HealthCheckResult{
				"database": {Status: "skip", Detail: "setup mode"},
			},
		}
	}

	ctx, cancel := context.WithTimeout(ctx, api.timeout)
	defer cancel()

	if api.checkDB == nil || api.checkDB(ctx) != nil {
		return http.StatusServiceUnavailable, dto.HealthResponse{
			Status: "not_ready",
			Mode:   api.mode,
			Checks: map[string]dto.HealthCheckResult{
				"database": {Status: "fail", Detail: "database ping failed"},
			},
		}
	}

	return http.StatusOK, dto.HealthResponse{
		Status: "ok",
		Mode:   api.mode,
		Checks: map[string]dto.HealthCheckResult{
			"database": {Status: "ok"},
		},
	}
}

func readyStateFromCode(code int) float64 {
	if code == http.StatusOK {
		return 1
	}
	return 0
}
