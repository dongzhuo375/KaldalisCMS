package v1

import (
	"KaldalisCMS/internal/api/middleware"
	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/core"
	repository "KaldalisCMS/internal/infra/repository/postgres"
	"KaldalisCMS/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MediaAPI struct {
	svc  *service.MediaService
	repo *repository.MediaRepository
}

func NewMediaAPI(svc *service.MediaService, repo *repository.MediaRepository) *MediaAPI {
	return &MediaAPI{svc: svc, repo: repo}
}

// RegisterRoutes registers media routes under /api/v1.
func (api *MediaAPI) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/media", api.Upload)
	rg.GET("/media", api.List)
	rg.DELETE("/media/:id", api.Delete)
	// per-post media library (references)
	rg.GET("/posts/:id/media", api.ListPostMedia)
}

func (api *MediaAPI) Upload(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
		return
	}

	asset, err := api.svc.CreateAssetFromUpload(c.Request.Context(), userID, file)
	if err != nil {
		switch err {
		case service.ErrUploadTooLarge:
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": err.Error()})
			return
		case service.ErrUnsupportedType:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"asset": dto.ToMediaAssetResponse(asset)})
}

func (api *MediaAPI) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	roleVal, _ := c.Get("kaldalis_user_role")
	role, _ := roleVal.(string)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	q := c.DefaultQuery("q", "")

	assets, total, err := api.svc.List(c.Request.Context(), role, userID, page, pageSize, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.MediaListResponse{Items: dto.ToMediaAssetResponses(assets), Total: total, Page: page, PageSize: pageSize})
}

func (api *MediaAPI) Delete(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	roleVal, _ := c.Get("kaldalis_user_role")
	role, _ := roleVal.(string)

	if err := api.svc.DeleteAs(c.Request.Context(), role, userID, uint(id64)); err != nil {
		switch err {
		case service.ErrAssetReferenced:
			cnt, _ := api.repo.CountReferences(c.Request.Context(), uint(id64))
			c.JSON(http.StatusConflict, gin.H{"error": err.Error(), "references": cnt})
			return
		case repository.ErrMediaNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		case core.ErrPermission:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (api *MediaAPI) ListPostMedia(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}
	purpose := c.Query("purpose")
	var p *string
	if purpose != "" {
		p = &purpose
	}

	assets, err := api.svc.ListPostMedia(c.Request.Context(), uint(id64), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": dto.ToMediaAssetResponses(assets)})
}
