package middleware

import (
	"net/http"
	"strings"

	infraauth "KaldalisCMS/internal/infra/auth"

	"github.com/gin-gonic/gin"
)

const csrfHeader = "X-CSRF-Token"

// CSRFMiddlewareWithManager 返回一个 Gin 中间件（double-submit cookie），使用 infra auth.Manager 的 CSRFCookie 名称
func CSRFMiddlewareWithManager(mgr *infraauth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		// 安全方法跳过 CSRF 校验
		if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
			c.Next()
			return
		}

		// 读取 cookie 与 header
		cookie, err := c.Request.Cookie(mgr.CSRFCookie)
		if err != nil || cookie == nil || cookie.Value == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing csrf cookie"})
			return
		}
		headerToken := c.GetHeader(csrfHeader)
		if headerToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing csrf header"})
			return
		}
		if headerToken != cookie.Value {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid csrf token"})
			return
		}

		// Optional: 简单校验 Origin（如果有）和 Host 是否匹配（增强防护）
		origin := c.GetHeader("Origin")
		if origin != "" {
			if !strings.Contains(origin, c.Request.Host) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid origin"})
				return
			}
		}

		c.Next()
	}
}
