package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() { gin.SetMode(gin.TestMode) }

func TestCORSMiddleware_HeadersOnGET(t *testing.T) {
	r := gin.New()
	r.Use(CORSMiddleware())
	r.GET("/x", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: %d", w.Code)
	}
	h := w.Header()
	if h.Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Fatalf("allow-origin: %q", h.Get("Access-Control-Allow-Origin"))
	}
	if h.Get("Access-Control-Allow-Credentials") != "true" {
		t.Fatal("credentials header missing")
	}
	if h.Get("Access-Control-Allow-Methods") == "" {
		t.Fatal("methods header missing")
	}
}

func TestCORSMiddleware_OPTIONSShortCircuit(t *testing.T) {
	handlerCalled := false
	r := gin.New()
	r.Use(CORSMiddleware())
	r.OPTIONS("/x", func(c *gin.Context) {
		handlerCalled = true
		c.String(http.StatusOK, "should not reach")
	})

	req := httptest.NewRequest(http.MethodOptions, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status: %d (want 204)", w.Code)
	}
	if handlerCalled {
		t.Fatal("downstream handler must not run on preflight")
	}
}
