package v1

import (
	"KaldalisCMS/internal/api/errorx"
	"KaldalisCMS/internal/api/middleware"
	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/core"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminPostAPI serves post-management endpoints under /api/v1/admin/posts.
// The path is shared by two back-office personas:
//   - regular users manage only their own drafts
//   - admin roles manage all posts and can execute publish/offline operations
//
// Ownership filtering and workflow authorization stay in the service layer so the HTTP layer
// only needs to extract actor context and translate errors into responses.
type AdminPostAPI struct {
	service core.PostService
}

func NewAdminPostAPI(service core.PostService) *AdminPostAPI {
	return &AdminPostAPI{service: service}
}

// GetPosts returns the management list for the current actor.
// Admins receive the full set, while regular users only receive their own drafts.
// @Summary List manageable posts
// @Description Returns admin list for current actor scope (own drafts or all posts).
// @Tags admin-posts
// @Produce json
// @Success 200 {array} dto.PostResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Failure 504 {object} dto.ErrorResponse
// @Security CookieAuth
// @Security CSRFToken
// @Router /admin/posts [get]
func (api *AdminPostAPI) GetPosts(c *gin.Context) {
	actorUserID, actorRole, ok := getPostActor(c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	posts, err := api.service.ListAdminPosts(ctx, actorUserID, actorRole)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			errorx.RespondTimeoutError(c, "list admin posts timed out")
			return
		}
		respondPostWorkflowError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dto.ToPostListResponse(posts))
}

// GetPostByID returns a single manageable post for the current actor.
// @Summary Get manageable post by ID
// @Description Returns one post visible to current actor scope in admin workflow.
// @Tags admin-posts
// @Produce json
// @Param id path int true "post id"
// @Success 200 {object} dto.PostResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 504 {object} dto.ErrorResponse
// @Security CookieAuth
// @Security CSRFToken
// @Router /admin/posts/{id} [get]
func (api *AdminPostAPI) GetPostByID(c *gin.Context) {
	id, ok := parsePostID(c)
	if !ok {
		return
	}

	actorUserID, actorRole, ok := getPostActor(c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	post, err := api.service.GetAdminPostByID(ctx, id, actorUserID, actorRole)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			errorx.RespondTimeoutError(c, "get admin post timed out")
			return
		}
		respondPostWorkflowError(c, err, http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, dto.ToPostResponse(&post))
}

// CreatePost creates a new management post.
// The service layer always persists it as Draft and binds ownership to the authenticated actor.
// @Summary Create post draft
// @Description Create a new draft post under admin workflow.
// @Tags admin-posts
// @Accept json
// @Produce json
// @Param body body dto.CreatePostRequest true "create post payload"
// @Success 201 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Failure 504 {object} dto.ErrorResponse
// @Security CookieAuth
// @Security CSRFToken
// @Router /admin/posts [post]
func (api *AdminPostAPI) CreatePost(c *gin.Context) {
	var createReq dto.CreatePostRequest
	if err := c.ShouldBindJSON(&createReq); err != nil {
		errorx.RespondValidationError(c, "invalid request body", map[string]any{"reason": err.Error()})
		return
	}

	actorUserID, actorRole, ok := getPostActor(c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := api.service.CreateAdminPost(ctx, actorUserID, actorRole, *createReq.ToEntity()); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			errorx.RespondTimeoutError(c, "create post timed out")
			return
		}
		respondPostWorkflowError(c, err, http.StatusInternalServerError)
		return
	}

	errorx.RespondMessage(c, http.StatusCreated, "post created successfully")
}

// UpdatePost updates editable post content fields.
// Publication status is intentionally excluded and managed by dedicated workflow endpoints.
// @Summary Update post draft fields
// @Description Update editable fields for one post in admin workflow.
// @Tags admin-posts
// @Accept json
// @Produce json
// @Param id path int true "post id"
// @Param body body dto.UpdatePostRequest true "update post payload"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 504 {object} dto.ErrorResponse
// @Security CookieAuth
// @Security CSRFToken
// @Router /admin/posts/{id} [put]
func (api *AdminPostAPI) UpdatePost(c *gin.Context) {
	id, ok := parsePostID(c)
	if !ok {
		return
	}

	actorUserID, actorRole, ok := getPostActor(c)
	if !ok {
		return
	}

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorx.RespondValidationError(c, "invalid request body", map[string]any{"reason": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := api.service.UpdateAdminPost(ctx, id, req.ToPatch(), actorUserID, actorRole); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			errorx.RespondTimeoutError(c, "update post timed out")
			return
		}
		respondPostWorkflowError(c, err, http.StatusNotFound)
		return
	}

	errorx.RespondMessage(c, http.StatusOK, "updated")
}

// PublishPost transitions a post from Draft to Published.
// @Summary Publish post
// @Description Transition a draft post to published status.
// @Tags admin-posts
// @Produce json
// @Param id path int true "post id"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 504 {object} dto.ErrorResponse
// @Security CookieAuth
// @Security CSRFToken
// @Router /admin/posts/{id}/publish [post]
func (api *AdminPostAPI) PublishPost(c *gin.Context) {
	id, ok := parsePostID(c)
	if !ok {
		return
	}

	_, actorRole, ok := getPostActor(c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := api.service.PublishAdminPost(ctx, id, actorRole); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			errorx.RespondTimeoutError(c, "publish post timed out")
			return
		}
		respondPostWorkflowError(c, err, http.StatusBadRequest)
		return
	}

	errorx.RespondMessage(c, http.StatusOK, "post published successfully")
}

// DraftPost performs the minimal offline action by moving a post back to Draft.
// @Summary Move post to draft
// @Description Transition a published post back to draft status.
// @Tags admin-posts
// @Produce json
// @Param id path int true "post id"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 504 {object} dto.ErrorResponse
// @Security CookieAuth
// @Security CSRFToken
// @Router /admin/posts/{id}/draft [post]
func (api *AdminPostAPI) DraftPost(c *gin.Context) {
	id, ok := parsePostID(c)
	if !ok {
		return
	}

	_, actorRole, ok := getPostActor(c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := api.service.MovePostToDraft(ctx, id, actorRole); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			errorx.RespondTimeoutError(c, "move post to draft timed out")
			return
		}
		respondPostWorkflowError(c, err, http.StatusBadRequest)
		return
	}

	errorx.RespondMessage(c, http.StatusOK, "post moved to draft successfully")
}

// DeletePost removes a post from the system.
// @Summary Delete post
// @Description Permanently delete one post under admin workflow authorization.
// @Tags admin-posts
// @Produce json
// @Param id path int true "post id"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 504 {object} dto.ErrorResponse
// @Security CookieAuth
// @Security CSRFToken
// @Router /admin/posts/{id} [delete]
func (api *AdminPostAPI) DeletePost(c *gin.Context) {
	id, ok := parsePostID(c)
	if !ok {
		return
	}

	_, actorRole, ok := getPostActor(c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := api.service.DeleteAdminPost(ctx, id, actorRole); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			errorx.RespondTimeoutError(c, "delete post timed out")
			return
		}
		respondPostWorkflowError(c, err, http.StatusNotFound)
		return
	}

	errorx.RespondMessage(c, http.StatusOK, "post deleted successfully")
}

func getPostActor(c *gin.Context) (uint, string, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		errorx.RespondError(c, http.StatusUnauthorized, core.CodeUnauthorized, "unauthorized", nil)
		return 0, "", false
	}

	role, ok := middleware.GetUserRole(c)
	if !ok {
		errorx.RespondError(c, http.StatusUnauthorized, core.CodeUnauthorized, "unauthorized", nil)
		return 0, "", false
	}

	return userID, role, true
}

func respondPostWorkflowError(c *gin.Context, err error, defaultStatus int) {
	status := defaultStatus
	if errors.Is(err, core.ErrInternalError) {
		status = http.StatusInternalServerError
	}
	errorx.RespondErrorByCore(c, err, status, nil)
}
