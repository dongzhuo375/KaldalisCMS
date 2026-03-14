//go:build !swagger

package router

import "github.com/gin-gonic/gin"

func registerSwaggerRoutes(_ *gin.Engine, _ SwaggerOptions) {
	// Swagger integration is excluded from non-swagger builds.
}
