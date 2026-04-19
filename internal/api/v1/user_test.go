package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	"KaldalisCMS/pkg/auth"

	"github.com/gin-gonic/gin"
)

// --- fakes ---

type fakeUserService struct {
	createFn    func(ctx context.Context, u entity.User) error
	verifyFn    func(ctx context.Context, username, password string) (entity.User, error)
	loginFn     func(ctx context.Context, username, password string) (entity.User, error)
	getByIDFn   func(ctx context.Context, id uint) (entity.User, error)
}

func (f *fakeUserService) CreateUser(ctx context.Context, u entity.User) error {
	return f.createFn(ctx, u)
}
func (f *fakeUserService) VerifyUser(ctx context.Context, username, password string) (entity.User, error) {
	return f.verifyFn(ctx, username, password)
}
func (f *fakeUserService) Login(ctx context.Context, username, password string) (entity.User, error) {
	return f.loginFn(ctx, username, password)
}
func (f *fakeUserService) GetUserByID(ctx context.Context, id uint) (entity.User, error) {
	return f.getByIDFn(ctx, id)
}
func (f *fakeUserService) Logout() {}

type fakeSessionManager struct {
	establishFn func(w http.ResponseWriter, userID uint, role string) error
	destroyFn   func(w http.ResponseWriter)
	authFn      func(r *http.Request) (*auth.CustomClaims, error)
	csrfFn      func(r *http.Request, hash string) error
	ttl         time.Duration
}

func (f *fakeSessionManager) EstablishSession(w http.ResponseWriter, uid uint, role string) error {
	return f.establishFn(w, uid, role)
}
func (f *fakeSessionManager) DestroySession(w http.ResponseWriter) {
	if f.destroyFn != nil {
		f.destroyFn(w)
	}
}
func (f *fakeSessionManager) Authenticate(r *http.Request) (*auth.CustomClaims, error) {
	return f.authFn(r)
}
func (f *fakeSessionManager) ValidateCSRF(r *http.Request, hash string) error {
	return f.csrfFn(r, hash)
}
func (f *fakeSessionManager) GetTTL() time.Duration { return f.ttl }

// --- helpers ---

func newUserRouter(svc core.UserService, sm core.SessionManager, actorMW gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	api := NewUserAPI(svc, sm)
	grp := r.Group("")
	if actorMW != nil {
		grp.Use(actorMW)
	}
	grp.POST("/users/register", api.Register)
	grp.POST("/users/login", api.Login)
	grp.GET("/users/profile", api.GetProfile)
	grp.POST("/users/logout", api.Logout)
	return r
}

// --- tests ---

func TestUserAPI_Register_InvalidBody(t *testing.T) {
	r := newUserRouter(&fakeUserService{}, nil, nil)
	// Missing required fields
	w := doJSON(r, http.MethodPost, "/users/register", map[string]any{
		"username": "ab", // min=3
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: %d body=%s", w.Code, w.Body.String())
	}
}

func TestUserAPI_Register_Success(t *testing.T) {
	var captured entity.User
	svc := &fakeUserService{
		createFn: func(ctx context.Context, u entity.User) error {
			captured = u
			return nil
		},
	}
	r := newUserRouter(svc, nil, nil)
	w := doJSON(r, http.MethodPost, "/users/register", dto.UserRegisterRequest{
		Username: "alice",
		Password: "secret123",
		Email:    "a@b.com",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("status: %d body=%s", w.Code, w.Body.String())
	}
	if captured.Role != "user" {
		t.Fatalf("default role not applied: %q", captured.Role)
	}
}

func TestUserAPI_Register_DuplicateUser(t *testing.T) {
	svc := &fakeUserService{
		createFn: func(ctx context.Context, u entity.User) error {
			return core.ErrDuplicate
		},
	}
	r := newUserRouter(svc, nil, nil)
	w := doJSON(r, http.MethodPost, "/users/register", dto.UserRegisterRequest{
		Username: "alice",
		Password: "secret123",
		Email:    "a@b.com",
	})
	if w.Code != http.StatusConflict {
		t.Fatalf("status: %d", w.Code)
	}
}

func TestUserAPI_Login_InvalidBody(t *testing.T) {
	r := newUserRouter(&fakeUserService{}, &fakeSessionManager{}, nil)
	w := doJSON(r, http.MethodPost, "/users/login", map[string]any{})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", w.Code)
	}
}

func TestUserAPI_Login_BadCredentials(t *testing.T) {
	svc := &fakeUserService{
		loginFn: func(ctx context.Context, username, password string) (entity.User, error) {
			return entity.User{}, core.ErrInvalidCredentials
		},
	}
	r := newUserRouter(svc, &fakeSessionManager{}, nil)
	w := doJSON(r, http.MethodPost, "/users/login", dto.UserLoginRequest{
		Username: "alice", Password: "wrong",
	})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: %d", w.Code)
	}
}

func TestUserAPI_Login_Success(t *testing.T) {
	svc := &fakeUserService{
		loginFn: func(ctx context.Context, username, password string) (entity.User, error) {
			return entity.User{ID: 1, Username: "alice", Email: "a@b.com", Role: "admin"}, nil
		},
	}
	sm := &fakeSessionManager{
		establishFn: func(w http.ResponseWriter, uid uint, role string) error { return nil },
		ttl:         24 * time.Hour,
	}
	r := newUserRouter(svc, sm, nil)
	w := doJSON(r, http.MethodPost, "/users/login", dto.UserLoginRequest{
		Username: "alice", Password: "pw",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("status: %d body=%s", w.Code, w.Body.String())
	}
	var got dto.LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.User.ID != 1 || got.User.Role != "admin" {
		t.Fatalf("user: %+v", got.User)
	}
	if got.ExpiresAt == "" {
		t.Fatal("expires_at missing")
	}
}

func TestUserAPI_Login_SessionEstablishFails(t *testing.T) {
	svc := &fakeUserService{
		loginFn: func(ctx context.Context, username, password string) (entity.User, error) {
			return entity.User{ID: 1}, nil
		},
	}
	sm := &fakeSessionManager{
		establishFn: func(w http.ResponseWriter, uid uint, role string) error {
			return core.ErrInternalError
		},
		ttl: time.Hour,
	}
	r := newUserRouter(svc, sm, nil)
	w := doJSON(r, http.MethodPost, "/users/login", dto.UserLoginRequest{
		Username: "alice", Password: "pw",
	})
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", w.Code)
	}
}

func TestUserAPI_GetProfile_Unauthenticated(t *testing.T) {
	r := newUserRouter(&fakeUserService{}, nil, nil)
	w := doRequest(r, http.MethodGet, "/users/profile")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: %d", w.Code)
	}
}

func TestUserAPI_GetProfile_Success(t *testing.T) {
	svc := &fakeUserService{
		getByIDFn: func(ctx context.Context, id uint) (entity.User, error) {
			return entity.User{ID: id, Username: "alice", Email: "a@b.com", Role: "admin"}, nil
		},
	}
	actor := injectActor(7, "admin")
	r := newUserRouter(svc, nil, actor)
	w := doRequest(r, http.MethodGet, "/users/profile")
	if w.Code != http.StatusOK {
		t.Fatalf("status: %d body=%s", w.Code, w.Body.String())
	}
	var got dto.LoginUserResponse
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if got.ID != 7 || got.Username != "alice" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestUserAPI_GetProfile_NotFound(t *testing.T) {
	svc := &fakeUserService{
		getByIDFn: func(ctx context.Context, id uint) (entity.User, error) {
			return entity.User{}, core.ErrNotFound
		},
	}
	r := newUserRouter(svc, nil, injectActor(7, "admin"))
	w := doRequest(r, http.MethodGet, "/users/profile")
	if w.Code != http.StatusNotFound {
		t.Fatalf("status: %d", w.Code)
	}
}

func TestUserAPI_Logout(t *testing.T) {
	destroyed := false
	sm := &fakeSessionManager{
		destroyFn: func(w http.ResponseWriter) { destroyed = true },
	}
	r := newUserRouter(&fakeUserService{}, sm, nil)
	req := httptest.NewRequest(http.MethodPost, "/users/logout", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status: %d", w.Code)
	}
	if !destroyed {
		t.Fatal("session not destroyed")
	}
}
