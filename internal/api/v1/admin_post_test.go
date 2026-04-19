package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"

	"github.com/gin-gonic/gin"
)

// These string literals mirror the unexported context keys in internal/api/middleware/auth.go.
// If the middleware renames them, these tests will fail — that is the intended contract fence.
const (
	testCtxUserIDKey   = "kaldalis_user_id"
	testCtxUserRoleKey = "kaldalis_user_role"
)

// injectActor simulates the auth middleware by pre-populating the actor identity.
// Pass uid=0 to simulate an unauthenticated request (key absent).
func injectActor(uid uint, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if uid != 0 {
			c.Set(testCtxUserIDKey, uid)
			c.Set(testCtxUserRoleKey, role)
		}
		c.Next()
	}
}

func newAdminRouter(svc core.PostService, actor gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	api := NewAdminPostAPI(svc)
	grp := r.Group("/admin/posts")
	grp.Use(actor)
	grp.GET("", api.GetPosts)
	grp.GET("/:id", api.GetPostByID)
	grp.POST("", api.CreatePost)
	grp.PUT("/:id", api.UpdatePost)
	grp.DELETE("/:id", api.DeletePost)
	grp.POST("/:id/publish", api.PublishPost)
	grp.POST("/:id/draft", api.DraftPost)
	return r
}

func doJSON(r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	var reader *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reader = bytes.NewReader(b)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAdminPostAPI_Unauthenticated(t *testing.T) {
	// No actor injected → handler must return 401 without reaching the service.
	svc := &fakePostService{} // any service call would nil-panic, proving the guard
	r := newAdminRouter(svc, injectActor(0, ""))

	w := doJSON(r, http.MethodGet, "/admin/posts", nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: %d body=%s", w.Code, w.Body.String())
	}
	var got dto.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if got.Code != string(core.CodeUnauthorized) {
		t.Fatalf("code: %s", got.Code)
	}
}

func TestAdminPostAPI_GetPosts_Success(t *testing.T) {
	svc := &fakePostService{
		listAdminFn: func(ctx context.Context, uid uint, role string) ([]entity.Post, error) {
			if uid != 9 || role != "admin" {
				t.Fatalf("actor not propagated: uid=%d role=%s", uid, role)
			}
			return []entity.Post{{ID: 1}, {ID: 2}}, nil
		},
	}
	r := newAdminRouter(svc, injectActor(9, "admin"))
	w := doJSON(r, http.MethodGet, "/admin/posts", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status: %d body=%s", w.Code, w.Body.String())
	}
	var got []dto.PostResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("len: %d", len(got))
	}
}

func TestAdminPostAPI_CreatePost_InvalidBody(t *testing.T) {
	svc := &fakePostService{} // service must not be called on binding failure
	r := newAdminRouter(svc, injectActor(9, "admin"))

	w := doJSON(r, http.MethodPost, "/admin/posts", map[string]any{
		// Title is required (binding:"required") → binding should reject
		"content": "body",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", w.Code)
	}
	var got dto.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if got.Code != string(core.CodeValidationFailed) {
		t.Fatalf("code: %s", got.Code)
	}
}

func TestAdminPostAPI_CreatePost_Success(t *testing.T) {
	var captured entity.Post
	svc := &fakePostService{
		createAdminFn: func(ctx context.Context, uid uint, role string, p entity.Post) (entity.Post, error) {
			captured = p
			p.ID = 77
			return p, nil
		},
	}
	r := newAdminRouter(svc, injectActor(3, "editor"))

	w := doJSON(r, http.MethodPost, "/admin/posts", dto.CreatePostRequest{
		Title:   "Hello",
		Content: "body",
		Tags:    []uint{1, 2},
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("status: %d body=%s", w.Code, w.Body.String())
	}
	if captured.Title != "Hello" {
		t.Fatalf("title not forwarded: %q", captured.Title)
	}
	if len(captured.Tags) != 2 || captured.Tags[0].ID != 1 {
		t.Fatalf("tags not mapped: %+v", captured.Tags)
	}
	var got dto.PostResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != 77 {
		t.Fatalf("id: %d", got.ID)
	}
}

func TestAdminPostAPI_PublishPost_ConflictMapping(t *testing.T) {
	// Service returns a wrapped ErrConflict → handler must map to 409.
	svc := &fakePostService{
		publishAdminFn: func(ctx context.Context, id uint, uid uint, role string) error {
			return fmt.Errorf("publish: %w", core.ErrConflict)
		},
	}
	r := newAdminRouter(svc, injectActor(9, "admin"))
	w := doJSON(r, http.MethodPost, "/admin/posts/1/publish", nil)
	if w.Code != http.StatusConflict {
		t.Fatalf("status: %d", w.Code)
	}
	var got dto.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if got.Code != string(core.CodeConflict) {
		t.Fatalf("code: %s", got.Code)
	}
}

func TestAdminPostAPI_DeletePost_ForbiddenMapping(t *testing.T) {
	svc := &fakePostService{
		deleteAdminFn: func(ctx context.Context, id uint, uid uint, role string) error {
			return core.ErrPermission
		},
	}
	r := newAdminRouter(svc, injectActor(9, "editor"))
	w := doJSON(r, http.MethodDelete, "/admin/posts/5", nil)
	if w.Code != http.StatusForbidden {
		t.Fatalf("status: %d", w.Code)
	}
}

func TestAdminPostAPI_DeletePost_Success(t *testing.T) {
	svc := &fakePostService{
		deleteAdminFn: func(ctx context.Context, id uint, uid uint, role string) error {
			if id != 5 {
				t.Fatalf("id not parsed: %d", id)
			}
			return nil
		},
	}
	r := newAdminRouter(svc, injectActor(9, "admin"))
	w := doJSON(r, http.MethodDelete, "/admin/posts/5", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status: %d", w.Code)
	}
	var got dto.MessageResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if !strings.Contains(got.Message, "deleted") {
		t.Fatalf("message: %q", got.Message)
	}
}

func TestAdminPostAPI_DraftPost_Success(t *testing.T) {
	svc := &fakePostService{
		moveToDraftAdminFn: func(ctx context.Context, id uint, uid uint, role string) error {
			return nil
		},
	}
	r := newAdminRouter(svc, injectActor(9, "admin"))
	w := doJSON(r, http.MethodPost, "/admin/posts/1/draft", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status: %d", w.Code)
	}
}

func TestAdminPostAPI_UpdatePost_PatchApplied(t *testing.T) {
	var capturedPatch entity.PostPatch
	svc := &fakePostService{
		updateAdminFn: func(ctx context.Context, id uint, patch entity.PostPatch, uid uint, role string) error {
			capturedPatch = patch
			return nil
		},
	}
	r := newAdminRouter(svc, injectActor(9, "admin"))

	w := doJSON(r, http.MethodPut, "/admin/posts/3", dto.UpdatePostRequest{
		Title: strPtr("new title"),
		Tags:  []uint{4, 5},
	})
	if w.Code != http.StatusOK {
		t.Fatalf("status: %d body=%s", w.Code, w.Body.String())
	}
	if capturedPatch.Title == nil || *capturedPatch.Title != "new title" {
		t.Fatalf("title patch lost: %+v", capturedPatch.Title)
	}
	if capturedPatch.Content != nil {
		t.Fatal("content patch should remain nil (not in request)")
	}
	if len(capturedPatch.Tags) != 2 {
		t.Fatalf("tags patch: %+v", capturedPatch.Tags)
	}
}

func TestAdminPostAPI_GetPostByID_NotFoundMapping(t *testing.T) {
	svc := &fakePostService{
		getAdminByIDFn: func(ctx context.Context, id uint, uid uint, role string) (entity.Post, error) {
			return entity.Post{}, core.ErrNotFound
		},
	}
	r := newAdminRouter(svc, injectActor(9, "admin"))
	w := doJSON(r, http.MethodGet, "/admin/posts/99", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status: %d", w.Code)
	}
}

func TestAdminPostAPI_InternalErrorMapping(t *testing.T) {
	// respondPostWorkflowError must upgrade to 500 when ErrInternalError is in the chain.
	svc := &fakePostService{
		getAdminByIDFn: func(ctx context.Context, id uint, uid uint, role string) (entity.Post, error) {
			return entity.Post{}, fmt.Errorf("db: %w", core.ErrInternalError)
		},
	}
	r := newAdminRouter(svc, injectActor(9, "admin"))
	w := doJSON(r, http.MethodGet, "/admin/posts/1", nil)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", w.Code)
	}
	// sanity: the wrapped error still satisfies errors.Is
	if !errors.Is(fmt.Errorf("db: %w", core.ErrInternalError), core.ErrInternalError) {
		t.Fatal("errors.Is broke — test premise invalid")
	}
}

func strPtr(s string) *string { return &s }
