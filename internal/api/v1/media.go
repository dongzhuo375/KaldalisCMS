package v1

import (
	"KaldalisCMS/internal/api/errorx"
	"KaldalisCMS/internal/api/middleware"
	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/core"
	repository "KaldalisCMS/internal/infra/repository/postgres"
	"KaldalisCMS/internal/service"
	"errors"
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

// Upload stores one media file owned by current user.
// @Summary Upload media asset
// @Description Upload one file and create a media asset record.
// @Tags media
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "media file"
// @Success 201 {object} dto.MediaUploadResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 413 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security CookieAuth
// @Security CSRFToken
// @Router /media [post]
func (api *MediaAPI) Upload(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		errorx.RespondError(c, http.StatusUnauthorized, core.CodeUnauthorized, "unauthorized", nil)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		errorx.RespondValidationError(c, "missing file", nil)
		return
	}

	asset, err := api.svc.CreateAssetFromUpload(c.Request.Context(), userID, file)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUploadTooLarge):
			errorx.RespondError(c, http.StatusRequestEntityTooLarge, core.CodeValidationFailed, "upload too large", nil)
			return
		case errors.Is(err, service.ErrUnsupportedType):
			errorx.RespondValidationError(c, "unsupported file type", nil)
			return
		case errors.Is(err, core.ErrInvalidInput):
			errorx.RespondValidationError(c, "invalid upload payload", map[string]any{"reason": err.Error()})
			return
		default:
			errorx.RespondErrorByCore(c, err, http.StatusInternalServerError, nil)
			return
		}
	}

	c.JSON(http.StatusCreated, dto.MediaUploadResponse{Asset: dto.ToMediaAssetResponse(asset)})
}

// List returns media assets visible to current actor.
// @Summary List media assets
// @Description List media assets for current user scope with pagination and query filter.
// @Tags media
// @Produce json
// @Param page query int false "page number" default(1)
// @Param page_size query int false "page size" default(20)
// @Param q query string false "search keyword"
// @Success 200 {object} dto.MediaListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security CookieAuth
// @Security CSRFToken
// @Router /media [get]
func (api *MediaAPI) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		errorx.RespondError(c, http.StatusUnauthorized, core.CodeUnauthorized, "unauthorized", nil)
		return
	}
	roleVal, _ := c.Get("kaldalis_user_role")
	role, _ := roleVal.(string)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	q := c.DefaultQuery("q", "")

	assets, total, err := api.svc.List(c.Request.Context(), role, userID, page, pageSize, q)
	if err != nil {
		errorx.RespondInternalError(c)
		return
	}
	c.JSON(http.StatusOK, dto.MediaListResponse{Items: dto.ToMediaAssetResponses(assets), Total: total, Page: page, PageSize: pageSize})
}

// Delete removes one media asset by id.
// @Summary Delete media asset
// @Description Delete one media asset if caller has permission and asset is not referenced.
// @Tags media
// @Produce json
// @Param id path int true "media asset id"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security CookieAuth
// @Security CSRFToken
// @Router /media/{id} [delete]
func (api *MediaAPI) Delete(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errorx.RespondValidationError(c, "invalid id", map[string]any{"id": c.Param("id")})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		errorx.RespondError(c, http.StatusUnauthorized, core.CodeUnauthorized, "unauthorized", nil)
		return
	}
	roleVal, _ := c.Get("kaldalis_user_role")
	role, _ := roleVal.(string)

	if err := api.svc.DeleteAs(c.Request.Context(), role, userID, uint(id64)); err != nil {
		switch {
		case errors.Is(err, service.ErrAssetReferenced):
			cnt, _ := api.repo.CountReferences(c.Request.Context(), uint(id64))
			errorx.RespondError(c, http.StatusConflict, core.CodeConflict, "asset is referenced", map[string]any{"references": cnt})
			return
		case errors.Is(err, core.ErrNotFound):
			errorx.RespondError(c, http.StatusNotFound, core.CodeNotFound, "resource not found", nil)
			return
		case errors.Is(err, core.ErrPermission):
			errorx.RespondError(c, http.StatusForbidden, core.CodeForbidden, "permission denied", nil)
			return
		default:
			errorx.RespondInternalError(c)
			return
		}
	}
	errorx.RespondMessage(c, http.StatusOK, "deleted")
}

// ListPostMedia lists assets referenced by one post.
// @Summary List post media references
// @Description List media assets referenced by one post, optionally filtered by purpose.
// @Tags media
// @Produce json
// @Param id path int true "post id"
// @Param purpose query string false "reference purpose: content|cover"
// @Success 200 {object} dto.MediaItemsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /posts/{id}/media [get]
func (api *MediaAPI) ListPostMedia(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errorx.RespondValidationError(c, "invalid post id", map[string]any{"id": c.Param("id")})
		return
	}
	purpose := c.Query("purpose")
	var p *string
	if purpose != "" {
		p = &purpose
	}

	assets, err := api.svc.ListPostMedia(c.Request.Context(), uint(id64), p)
	if err != nil {
		errorx.RespondInternalError(c)
		return
	}
	c.JSON(http.StatusOK, dto.MediaItemsResponse{Items: dto.ToMediaAssetResponses(assets)})
}
