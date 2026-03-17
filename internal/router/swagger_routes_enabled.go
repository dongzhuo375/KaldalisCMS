//go:build swagger

package router

import (
	"KaldalisCMS/internal/core"
	docs "KaldalisCMS/internal/docs"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	openAPIv3Once sync.Once
	openAPIv3JSON []byte
	openAPIv3Err  error
)

func registerSwaggerRoutes(r *gin.Engine, opts SwaggerOptions) {
	if !opts.Enabled {
		return
	}

	if opts.Title != "" {
		docs.SwaggerInfo.Title = opts.Title
	}
	if opts.Version != "" {
		docs.SwaggerInfo.Version = opts.Version
	}
	if opts.Description != "" {
		docs.SwaggerInfo.Description = opts.Description
	}

	swaggerPath := normalizeSwaggerPath(opts.Path)
	openAPI3Path := swaggerPath + "-openapi3.json"

	r.GET(openAPI3Path, func(c *gin.Context) {
		spec, err := getOpenAPI3Spec()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    string(core.CodeInternalError),
				"message": "failed to build openapi spec",
				"details": map[string]any{"reason": err.Error()},
			})
			return
		}
		c.Data(http.StatusOK, "application/json; charset=utf-8", spec)
	})

	r.GET(swaggerPath+"/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(openAPI3Path)))
}

func getOpenAPI3Spec() ([]byte, error) {
	openAPIv3Once.Do(func() {
		var docV2 openapi2.T
		if err := json.Unmarshal([]byte(docs.SwaggerInfo.ReadDoc()), &docV2); err != nil {
			openAPIv3Err = err
			return
		}
		docV3, err := openapi2conv.ToV3(&docV2)
		if err != nil {
			openAPIv3Err = err
			return
		}
		openAPIv3JSON, openAPIv3Err = json.Marshal(docV3)
	})
	return openAPIv3JSON, openAPIv3Err
}
