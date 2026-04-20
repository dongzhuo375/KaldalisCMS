package middleware

import (
	"KaldalisCMS/internal/api/errorx"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

var httpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "kaldalis",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total HTTP requests grouped by method, route, status and api error code.",
	},
	[]string{"method", "route", "status", "code"},
)

var httpRequestDurationSeconds = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "kaldalis",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "HTTP request latency in seconds grouped by method and route.",
		Buckets:   prometheus.DefBuckets,
	},
	[]string{"method", "route"},
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDurationSeconds)
}

// RequestContext ensures every request has a stable request_id in context and headers.
func RequestContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(errorx.HeaderRequestID)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Set(errorx.CtxRequestIDKey, requestID)
		c.Writer.Header().Set(errorx.HeaderRequestID, requestID)
		c.Next()
	}
}

// RecoverAsContract converts panic paths into the same error envelope used by handlers.
func RecoverAsContract() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				reqID, _ := c.Get(errorx.CtxRequestIDKey)
				log.Printf("level=error event=panic request_id=%v method=%s path=%s panic=%v", reqID, c.Request.Method, c.Request.URL.Path, recovered)
				errorx.AbortInternalError(c)
			}
		}()
		c.Next()
	}
}

// ObserveHTTP emits structured logs and Prometheus metrics for all requests.
func ObserveHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		started := time.Now()
		c.Next()

		status := c.Writer.Status()
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		code := "-"
		if val, ok := c.Get(errorx.CtxErrorCodeKey); ok {
			if s, ok := val.(string); ok && s != "" {
				code = s
			}
		}
		httpRequestsTotal.WithLabelValues(c.Request.Method, route, fmt.Sprintf("%d", status), code).Inc()
		httpRequestDurationSeconds.WithLabelValues(c.Request.Method, route).Observe(time.Since(started).Seconds())

		reqID, _ := c.Get(errorx.CtxRequestIDKey)
		log.Printf("level=info event=http_request request_id=%v method=%s route=%s status=%d duration_ms=%d code=%s ip=%s", reqID, c.Request.Method, route, status, time.Since(started).Milliseconds(), code, c.ClientIP())
	}
}
