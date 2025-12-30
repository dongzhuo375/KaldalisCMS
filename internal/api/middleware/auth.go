package middleware

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/infra/auth"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ctxUserIDKey   = "kaldalis_user_id"
	ctxCsrfHashKey = "kaldalis_csrf_h"
)

// 识别
func OptionalAuth(sm core.SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := sm.Authenticate(c.Request)
		if err != nil {
			// 解析失败（比如过期或伪造），只要不是 "no token" 错误，就清理 Cookie 后放行
			if !errors.Is(err, auth.ErrNoToken) {
				sm.DestroySession(c.Writer)
			}
			c.Next()
			return
		}

		// 验证成功，将结果注入 Context
		c.Set(ctxUserIDKey, claims.UserID)
		c.Set(ctxCsrfHashKey, claims.CsrfH)
		c.Next()
	}
}

// 拦截
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不再解析 JWT
		if _, exists := c.Get(ctxUserIDKey); !exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}

func CSRFCheck(sm core.SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 放行审查
		//switch c.Request.Method {
		//case "GET", "HEAD", "OPTIONS", "TRACE":
		//	c.Next()
		//	return
		//}

		// 只有登录用户才检查 (从 Context 拿指纹)
		val, exists := c.Get(ctxCsrfHashKey)
		if !exists {
			c.Next()
			return
		}

		// CSRF校验
		if err := sm.ValidateCSRF(c.Request, val.(string)); err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.Next()
	}
}

func GetUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get(ctxUserIDKey)
	if !exists {
		return 0, false
	}

	uid, ok := val.(uint)
	return uid, ok
}
