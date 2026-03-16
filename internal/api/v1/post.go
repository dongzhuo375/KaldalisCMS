package v1

import (
	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/core"
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// PublicPostAPI serves anonymous/public article endpoints.
// It is intentionally restricted to published content so visibility rules remain
// stable regardless of caller identity.
type PublicPostAPI struct {
	service core.PostService
}

func NewPublicPostAPI(service core.PostService) *PublicPostAPI {
	return &PublicPostAPI{service: service}
}

func parsePostID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		respondValidationError(c, "invalid post id", map[string]any{"id": c.Param("id")})
		return 0, false
	}
	return uint(id), true
}

// GetPosts returns only published posts for public consumers.
// @Summary List published posts
// @Description Public read-only endpoint for published content.
// @Tags posts
// @Produce json
// @Success 200 {array} dto.PostResponse
// @Failure 500 {object} dto.ErrorResponse
// @Failure 504 {object} dto.ErrorResponse
// @Router /posts [get]
func (api *PublicPostAPI) GetPosts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	posts, err := api.service.ListPublicPosts(ctx)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			respondTimeoutError(c, "list posts timed out")
			return
		}
		respondInternalError(c)
		return
	}

	c.JSON(http.StatusOK, dto.ToPostListResponse(posts))
}

// GetPostByID returns a single published post.
// Drafts are intentionally invisible on this endpoint to avoid leaking unpublished content.
// @Summary Get published post
// @Description Public endpoint that returns one published post by numeric ID.
// @Tags posts
// @Produce json
// @Param id path int true "post id"
// @Success 200 {object} dto.PostResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 504 {object} dto.ErrorResponse
// @Router /posts/{id} [get]
func (api *PublicPostAPI) GetPostByID(c *gin.Context) {
	id, ok := parsePostID(c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	post, err := api.service.GetPublicPostByID(ctx, id)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			respondTimeoutError(c, "get post timed out")
			return
		}
		respondErrorByCore(c, err, http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, dto.ToPostResponse(&post))
}
