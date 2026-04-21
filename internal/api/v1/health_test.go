package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"KaldalisCMS/internal/api/v1/dto"

	"github.com/gin-gonic/gin"
)

func newHealthRouter(api *HealthAPI) *gin.Engine {
	r := gin.New()
	api.RegisterRootRoutes(r)
	return r
}

func TestHealthAPI_Healthz_AlwaysOK(t *testing.T) {
	api := &HealthAPI{mode: "app"}
	w := doRequest(newHealthRouter(api), http.MethodGet, "/healthz")
	if w.Code != http.StatusOK {
		t.Fatalf("status: %d", w.Code)
	}
	var got dto.HealthResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if got.Status != "ok" || got.Mode != "app" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestHealthAPI_Readyz_SetupMode_503(t *testing.T) {
	api := NewSetupHealthAPI()
	w := doRequest(newHealthRouter(api), http.MethodGet, "/readyz")
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status: %d", w.Code)
	}
	var got dto.HealthResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if got.Status != "not_ready" || got.Mode != "setup" {
		t.Fatalf("unexpected: %+v", got)
	}
	dbCheck, ok := got.Checks["database"]
	if !ok || dbCheck.Status != "skip" {
		t.Fatalf("db check: %+v", got.Checks)
	}
}

func TestHealthAPI_Readyz_AppMode_DBHealthy(t *testing.T) {
	api := &HealthAPI{
		mode:    "app",
		checkDB: func(ctx context.Context) error { return nil },
	}
	w := doRequest(newHealthRouter(api), http.MethodGet, "/readyz")
	if w.Code != http.StatusOK {
		t.Fatalf("status: %d", w.Code)
	}
	var got dto.HealthResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if got.Status != "ok" {
		t.Fatalf("status: %s", got.Status)
	}
	if got.Checks["database"].Status != "ok" {
		t.Fatalf("db check: %+v", got.Checks)
	}
}

func TestHealthAPI_Readyz_AppMode_DBDown(t *testing.T) {
	api := &HealthAPI{
		mode:    "app",
		checkDB: func(ctx context.Context) error { return errors.New("conn refused") },
	}
	w := doRequest(newHealthRouter(api), http.MethodGet, "/readyz")
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status: %d", w.Code)
	}
	var got dto.HealthResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if got.Status != "not_ready" {
		t.Fatalf("status: %s", got.Status)
	}
	if got.Checks["database"].Status != "fail" {
		t.Fatalf("db check: %+v", got.Checks)
	}
}

func TestHealthAPI_Readyz_AppMode_NilCheckDB(t *testing.T) {
	api := &HealthAPI{mode: "app", checkDB: nil}
	w := doRequest(newHealthRouter(api), http.MethodGet, "/readyz")
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("nil checkDB should return 503, got %d", w.Code)
	}
}

func TestHealthAPI_Readyz_CacheHit(t *testing.T) {
	callCount := 0
	api := &HealthAPI{
		mode: "app",
		checkDB: func(ctx context.Context) error {
			callCount++
			return nil
		},
		readySuccessTTL: 10_000_000_000, // 10s — ensures cache won't expire during test
	}
	r := newHealthRouter(api)

	w1 := doRequest(r, http.MethodGet, "/readyz")
	if w1.Code != http.StatusOK {
		t.Fatalf("first call: %d", w1.Code)
	}
	if callCount != 1 {
		t.Fatalf("checkDB should have been called once, got %d", callCount)
	}

	w2 := doRequest(r, http.MethodGet, "/readyz")
	if w2.Code != http.StatusOK {
		t.Fatalf("second call: %d", w2.Code)
	}
	if callCount != 1 {
		t.Fatalf("second call should hit cache, checkDB called %d times", callCount)
	}
}
