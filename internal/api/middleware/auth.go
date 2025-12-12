package middleware

import (
	"net/http"
	"strconv"
	"strings"

	infraauth "KaldalisCMS/internal/infra/auth"

	"github.com/gin-gonic/gin"
)

const ctxUserIDKey = "kaldalis_user_id"

// RequireAuthWithManager 返回一个 Gin 中间件函数，使用 infra auth.Manager 来解析 token
func RequireAuthWithManager(mgr *infraauth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		if cookie, err := c.Request.Cookie(mgr.AuthCookie); err == nil && cookie != nil && cookie.Value != "" {
			tokenStr = cookie.Value
		} else {
			// 回退到 Authorization header（Bearer ...）
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		claims, err := mgr.Parse(tokenStr)
		if err != nil {
			// 清 cookie
			mgr.Logout(c.Writer)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if uidVal, exists := claims["userID"]; exists {
			switch v := uidVal.(type) {
			case float64:
				c.Set(ctxUserIDKey, int(v))
			case int:
				c.Set(ctxUserIDKey, v)
			case string:
				if i, err := strconv.Atoi(v); err == nil {
					c.Set(ctxUserIDKey, i)
				} else {
					c.Set(ctxUserIDKey, v)
				}
			}
		}

		c.Next()
	}
}

// OptionalAuthWithManager 如果有合法 token 则注入 user，否则不拦截
func OptionalAuthWithManager(mgr *infraauth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		if cookie, err := c.Request.Cookie(mgr.AuthCookie); err == nil && cookie != nil && cookie.Value != "" {
			tokenStr = cookie.Value
		} else {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}
		if tokenStr == "" {
			c.Next()
			return
		}
		claims, err := mgr.Parse(tokenStr)
		if err != nil {
			// 清 cookie（可选）
			mgr.Logout(c.Writer)
			c.Next()
			return
		}
		if uidVal, exists := claims["userID"]; exists {
			switch v := uidVal.(type) {
			case float64:
				c.Set(ctxUserIDKey, int(v))
			case int:
				c.Set(ctxUserIDKey, v)
			case string:
				if i, err := strconv.Atoi(v); err == nil {
					c.Set(ctxUserIDKey, i)
				} else {
					c.Set(ctxUserIDKey, v)
				}
			}
		}
		c.Next()
	}
}

// GetUserID 从 gin.Context 获取当前请求的 user id（int 或 string）
func GetUserID(c *gin.Context) (interface{}, bool) {
	v, ok := c.Get(ctxUserIDKey)
	return v, ok
}
