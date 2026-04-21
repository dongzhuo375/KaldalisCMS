package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/gin-gonic/gin"
)

// minimal Casbin model: role-based matcher.
const testModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`

func newTestEnforcer(t *testing.T, policies [][]string) *casbin.Enforcer {
	t.Helper()
	m, err := model.NewModelFromString(testModel)
	if err != nil {
		t.Fatal(err)
	}
	e, err := casbin.NewEnforcer(m)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range policies {
		if _, err := e.AddPolicy(p); err != nil {
			t.Fatal(err)
		}
	}
	return e
}

func TestAuthorize_Allowed(t *testing.T) {
	e := newTestEnforcer(t, [][]string{{"admin", "/x", "GET"}})

	reached := false
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set(ctxUserRoleKey, "admin"); c.Next() })
	r.Use(Authorize(e))
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

func TestAuthorize_Forbidden(t *testing.T) {
	e := newTestEnforcer(t, [][]string{{"admin", "/x", "GET"}})

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set(ctxUserRoleKey, "user"); c.Next() })
	r.Use(Authorize(e))
	r.GET("/x", func(c *gin.Context) { t.Fatal("must not reach") })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status: %d", w.Code)
	}
}

func TestAuthorize_AnonymousFallback(t *testing.T) {
	// No role in context → Authorize should treat as "anonymous".
	e := newTestEnforcer(t, [][]string{{"anonymous", "/public", "GET"}})

	r := gin.New()
	r.Use(Authorize(e))
	r.GET("/public", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/public", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("anonymous policy should allow /public: %d", w.Code)
	}
}

func TestAuthorize_AnonymousDeniedWhenNoPolicy(t *testing.T) {
	e := newTestEnforcer(t, nil)

	r := gin.New()
	r.Use(Authorize(e))
	r.GET("/x", func(c *gin.Context) { t.Fatal("must not reach") })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status: %d", w.Code)
	}
}
