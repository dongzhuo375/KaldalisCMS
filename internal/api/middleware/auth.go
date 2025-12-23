package middleware

import (
	"KaldalisCMS/internal/infra/auth"
	pkgauth "KaldalisCMS/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

const ctxUserIDKey = "kaldalis_user_id"

func OptionalAuth(cfg auth.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie(cfg.AuthCookie)
		if token == "" {
			ah := c.GetHeader("Authorization")
			const prefix = "Bearer "
			if len(ah) > len(prefix) && ah[:len(prefix)] == prefix {
				token = ah[len(prefix):]
			}
		}

		if token == "" {
			c.Next()
			return
		}

		claims, err := pkgauth.Parse(token, cfg.Secret)
		if err != nil {
			auth.DestroySession(c.Writer, cfg)
			c.Next()
			return
		}

		c.Set(ctxUserIDKey, claims.UserID)
		c.Next()
	}
}

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

func CSRFCheck(cfg auth.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 放行审查
		//switch c.Request.Method {
		//case "GET", "HEAD", "OPTIONS", "TRACE":
		//	c.Next()
		//	return
		//}

		//冗余，做后续全员csrf验证的扩展
		//_, loggedIn := c.Get(ctxUserIDKey)
		//if !loggedIn {
		//	c.Next()
		//	return
		//}

		// 双重提交匹配
		cookieVal, _ := c.Cookie(cfg.CSRFCookie)
		headerVal := c.GetHeader("X-CSRF-Token")

		if cookieVal == "" || headerVal != cookieVal {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "CSRF token mismatch or missing",
			})
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
