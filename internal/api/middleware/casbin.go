package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// Authorize uses Casbin to check if the current user has permission.
// It reads the user's role from the Gin context, which should have been set by an upstream authentication middleware.
func Authorize(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取主体 (Subject): 从 Gin context 中获取用户角色
		var sub string
		roleVal, exists := c.Get(ctxUserRoleKey) // ctxUserRoleKey is defined in auth.go
		if exists {
			sub = roleVal.(string)
		} else {
			// 如果 context 中没有角色 (例如，未登录或 token 无效)，则视为匿名用户
			sub = "anonymous"
		}

		// 2. 获取对象 (Object): 请求路径 (使用 c.FullPath() 来适配 RESTful 的 :id)
		obj := c.FullPath()

		// 3. 获取行为 (Action): HTTP 请求方法
		act := c.Request.Method

		// 4. 使用 Casbin 执行器进行权限检查
		ok, err := e.Enforce(sub, obj, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "authorization check error"})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		// 权限检查通过
		c.Next()
	}
}
