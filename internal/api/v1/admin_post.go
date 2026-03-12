package v1

import (
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
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Get admin posts timed out"})
			return
		}
		if errors.Is(err, core.ErrPermission) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ToPostListResponse(posts))
}

// GetPostByID returns a single manageable post for the current actor.
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
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Get admin post timed out"})
			return
		}
		if errors.Is(err, core.ErrPermission) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ToPostResponse(&post))
}

// CreatePost creates a new management post.
// The service layer always persists it as Draft and binds ownership to the authenticated actor.
func (api *AdminPostAPI) CreatePost(c *gin.Context) {
	var createReq dto.CreatePostRequest
	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Create post timed out"})
			return
		}
		respondPostWorkflowError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
}

// UpdatePost updates editable post content fields.
// Publication status is intentionally excluded and managed by dedicated workflow endpoints.
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := api.service.UpdateAdminPost(ctx, id, req.ToPatch(), actorUserID, actorRole); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Update post timed out"})
			return
		}
		respondPostWorkflowError(c, err, http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// PublishPost transitions a post from Draft to Published.
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
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Publish post timed out"})
			return
		}
		respondPostWorkflowError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post published successfully"})
}

// DraftPost performs the minimal offline action by moving a post back to Draft.
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
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Unpublish post timed out"})
			return
		}
		respondPostWorkflowError(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post moved to draft successfully"})
}

// DeletePost removes a post from the system.
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
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Delete post timed out"})
			return
		}
		respondPostWorkflowError(c, err, http.StatusNotFound)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func getPostActor(c *gin.Context) (uint, string, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user ID not found in context"})
		return 0, "", false
	}

	role, ok := middleware.GetUserRole(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user role not found in context"})
		return 0, "", false
	}

	return userID, role, true
}

func respondPostWorkflowError(c *gin.Context, err error, defaultStatus int) {
	status := defaultStatus
	switch {
	case errors.Is(err, core.ErrPermission):
		status = http.StatusForbidden
	case errors.Is(err, core.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, core.ErrInvalidInput):
		status = http.StatusBadRequest
	}
	c.JSON(status, gin.H{"error": err.Error()})
}
