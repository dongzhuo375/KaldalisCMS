package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"KaldalisCMS/internal/infra/auth"
	pkgauth "KaldalisCMS/pkg/auth"

	"github.com/gin-gonic/gin"
)

type fakeSession struct {
	authFn    func(r *http.Request) (*pkgauth.CustomClaims, error)
	destroyed bool
	csrfFn    func(r *http.Request, hash string) error
}

func (f *fakeSession) EstablishSession(w http.ResponseWriter, uid uint, role string) error {
	return nil
}
func (f *fakeSession) DestroySession(w http.ResponseWriter) { f.destroyed = true }
func (f *fakeSession) Authenticate(r *http.Request) (*pkgauth.CustomClaims, error) {
	return f.authFn(r)
}
func (f *fakeSession) ValidateCSRF(r *http.Request, hash string) error {
	return f.csrfFn(r, hash)
}
func (f *fakeSession) GetTTL() time.Duration { return time.Hour }

// --- OptionalAuth ---

func TestOptionalAuth_InjectsClaims(t *testing.T) {
	sm := &fakeSession{authFn: func(r *http.Request) (*pkgauth.CustomClaims, error) {
		return &pkgauth.CustomClaims{UserID: 7, Role: "admin", CsrfH: "hh"}, nil
	}}
	var gotUID any
	var gotRole any
	var gotCsrf any
	r := gin.New()
	r.Use(OptionalAuth(sm))
	r.GET("/x", func(c *gin.Context) {
		gotUID, _ = c.Get(ctxUserIDKey)
		gotRole, _ = c.Get(ctxUserRoleKey)
		gotCsrf, _ = c.Get(ctxCsrfHashKey)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if gotUID != uint(7) {
		t.Fatalf("uid: %v (%T)", gotUID, gotUID)
	}
	if gotRole != "admin" {
		t.Fatalf("role: %v", gotRole)
	}
	if gotCsrf != "hh" {
		t.Fatalf("csrf: %v", gotCsrf)
	}
}

func TestOptionalAuth_NoTokenPassesThrough(t *testing.T) {
	sm := &fakeSession{authFn: func(r *http.Request) (*pkgauth.CustomClaims, error) {
		return nil, auth.ErrNoToken
	}}
	r := gin.New()
	r.Use(OptionalAuth(sm))
	r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("should pass through: %d", w.Code)
	}
	if sm.destroyed {
		t.Fatal("no-token must not destroy session")
	}
}

func TestOptionalAuth_InvalidTokenDestroysSession(t *testing.T) {
	sm := &fakeSession{authFn: func(r *http.Request) (*pkgauth.CustomClaims, error) {
		return nil, errors.New("bad token")
	}}
	r := gin.New()
	r.Use(OptionalAuth(sm))
	r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: %d", w.Code)
	}
	if !sm.destroyed {
		t.Fatal("invalid token should trigger DestroySession")
	}
}

// --- RequireAuth ---

func TestRequireAuth_Missing(t *testing.T) {
	r := gin.New()
	r.Use(RequireAuth())
	r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: %d", w.Code)
	}
}

func TestRequireAuth_Present(t *testing.T) {
	reached := false
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set(ctxUserIDKey, uint(1)); c.Next() })
	r.Use(RequireAuth())
	r.GET("/x", func(c *gin.Context) {
		reached = true
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if !reached || w.Code != http.StatusOK {
		t.Fatalf("not reached: code=%d", w.Code)
	}
}

// --- CSRFCheck ---

func TestCSRFCheck_SafeMethodsSkip(t *testing.T) {
	sm := &fakeSession{csrfFn: func(r *http.Request, h string) error {
		t.Fatal("should not be called for GET")
		return nil
	}}
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set(ctxCsrfHashKey, "h"); c.Next() })
	r.Use(CSRFCheck(sm))
	r.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: %d", w.Code)
	}
}

func TestCSRFCheck_AnonymousSkip(t *testing.T) {
	// No CsrfH in context (anonymous) → skip validation.
	called := false
	sm := &fakeSession{csrfFn: func(r *http.Request, h string) error {
		called = true
		return nil
	}}
	r := gin.New()
	r.Use(CSRFCheck(sm))
	r.POST("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodPost, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: %d", w.Code)
	}
	if called {
		t.Fatal("validator must not run when no CsrfH in context")
	}
}

func TestCSRFCheck_MutationSuccess(t *testing.T) {
	sm := &fakeSession{csrfFn: func(r *http.Request, h string) error {
		if h != "expected-h" {
			t.Fatalf("hash not forwarded: %q", h)
		}
		return nil
	}}
	reached := false
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set(ctxCsrfHashKey, "expected-h"); c.Next() })
	r.Use(CSRFCheck(sm))
	r.POST("/x", func(c *gin.Context) { reached = true; c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodPost, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK || !reached {
		t.Fatalf("status: %d reached=%v", w.Code, reached)
	}
}

func TestCSRFCheck_Forbidden(t *testing.T) {
	sm := &fakeSession{csrfFn: func(r *http.Request, h string) error {
		return errors.New("mismatch")
	}}
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set(ctxCsrfHashKey, "h"); c.Next() })
	r.Use(CSRFCheck(sm))
	r.POST("/x", func(c *gin.Context) {
		t.Fatal("must not reach")
	})

	req := httptest.NewRequest(http.MethodPost, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status: %d", w.Code)
	}
}

// --- GetUserID / GetUserRole ---

func TestGetUserID_Missing(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	if _, ok := GetUserID(c); ok {
		t.Fatal("missing key should return false")
	}
}

func TestGetUserID_WrongType(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(ctxUserIDKey, "not-a-uint")
	if _, ok := GetUserID(c); ok {
		t.Fatal("non-uint type should not pass type assertion")
	}
}

func TestGetUserID_Present(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(ctxUserIDKey, uint(42))
	got, ok := GetUserID(c)
	if !ok || got != 42 {
		t.Fatalf("got=%d ok=%v", got, ok)
	}
}

func TestGetUserRole_Cases(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	if _, ok := GetUserRole(c); ok {
		t.Fatal("missing should be false")
	}
	c.Set(ctxUserRoleKey, "admin")
	if got, ok := GetUserRole(c); !ok || got != "admin" {
		t.Fatalf("got=%q ok=%v", got, ok)
	}
}
