package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newPublicRouter(svc core.PostService) *gin.Engine {
	r := gin.New()
	api := NewPublicPostAPI(svc)
	r.GET("/posts", api.GetPosts)
	r.GET("/posts/:id", api.GetPostByID)
	return r
}

func doRequest(r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestPublicPostAPI_GetPosts_Success(t *testing.T) {
	svc := &fakePostService{
		listPublicFn: func(ctx context.Context) ([]entity.Post, error) {
			return []entity.Post{{ID: 1, Title: "hello"}}, nil
		},
	}
	w := doRequest(newPublicRouter(svc), http.MethodGet, "/posts")

	if w.Code != http.StatusOK {
		t.Fatalf("status: %d body=%s", w.Code, w.Body.String())
	}
	var got []dto.PostResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != 1 {
		t.Fatalf("unexpected body: %+v", got)
	}
}

func TestPublicPostAPI_GetPosts_ServiceError(t *testing.T) {
	svc := &fakePostService{
		listPublicFn: func(ctx context.Context) ([]entity.Post, error) {
			return nil, core.ErrInternalError
		},
	}
	w := doRequest(newPublicRouter(svc), http.MethodGet, "/posts")
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", w.Code)
	}
}

func TestPublicPostAPI_GetPostByID_InvalidID(t *testing.T) {
	w := doRequest(newPublicRouter(&fakePostService{}), http.MethodGet, "/posts/abc")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", w.Code)
	}
	var got dto.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Code != string(core.CodeValidationFailed) {
		t.Fatalf("code: %s", got.Code)
	}
	if got.Details["id"] != "abc" {
		t.Fatalf("details: %+v", got.Details)
	}
}

func TestPublicPostAPI_GetPostByID_NotFound(t *testing.T) {
	svc := &fakePostService{
		getPublicByIDFn: func(ctx context.Context, id uint) (entity.Post, error) {
			return entity.Post{}, core.ErrNotFound
		},
	}
	w := doRequest(newPublicRouter(svc), http.MethodGet, "/posts/9")
	if w.Code != http.StatusNotFound {
		t.Fatalf("status: %d", w.Code)
	}
	var got dto.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Code != string(core.CodeNotFound) {
		t.Fatalf("code: %s", got.Code)
	}
}

func TestPublicPostAPI_GetPostByID_Success(t *testing.T) {
	svc := &fakePostService{
		getPublicByIDFn: func(ctx context.Context, id uint) (entity.Post, error) {
			return entity.Post{ID: id, Title: "hi"}, nil
		},
	}
	w := doRequest(newPublicRouter(svc), http.MethodGet, "/posts/42")
	if w.Code != http.StatusOK {
		t.Fatalf("status: %d", w.Code)
	}
	var got dto.PostResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != 42 || got.Title != "hi" {
		t.Fatalf("body: %+v", got)
	}
}
