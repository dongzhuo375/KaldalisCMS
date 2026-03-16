package v1

import (
	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserAPI struct {
	service core.UserService
	sm      core.SessionManager
}

func NewUserAPI(service core.UserService, sessionMgr core.SessionManager) *UserAPI {
	return &UserAPI{
		service: service,
		sm:      sessionMgr,
	}
}

// RegisterRoutes registers the user-related routes to the Gin router.
func (api *UserAPI) RegisterRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/register", api.Register)
		userRoutes.POST("/login", api.Login)
	}
}

// Register handles new user registration.
// @Summary Register user
// @Description Create a normal user account.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.UserRegisterRequest true "register payload"
// @Success 201 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/register [post]
func (api *UserAPI) Register(c *gin.Context) {
	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, "invalid request body", map[string]any{"reason": err.Error()})
		return
	}

	newUser := entity.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     "user", // Assign default role
	}

	ctx := c.Request.Context()
	if err := api.service.CreateUser(ctx, newUser); err != nil {
		respondErrorByCore(c, err, http.StatusInternalServerError, nil)
		return
	}

	respondMessage(c, http.StatusCreated, "user created successfully")
}

// Login authenticates user credentials and creates session cookies.
// @Summary Login
// @Description Authenticate and establish a cookie-based session.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.UserLoginRequest true "login payload"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/login [post]
func (a *UserAPI) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, "invalid request body", map[string]any{"reason": err.Error()})
		return
	}
	user, err := a.service.Login(ctx, req.Username, req.Password)
	if err != nil {
		respondErrorByCore(c, err, http.StatusUnauthorized, nil)
		return
	}
	if err := a.sm.EstablishSession(c.Writer, user.ID, user.Role); err != nil {
		respondInternalError(c)
		return
	}

	expiresAt := time.Now().Add(a.sm.GetTTL())

	c.JSON(http.StatusOK, dto.LoginResponse{
		Message: "Login successful",
		User: dto.LoginUserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
		ExpiresAt: expiresAt.Format(time.RFC3339),
	})
}

// Logout clears the current cookie-based session.
// @Summary Logout
// @Description Destroy the current session cookies.
// @Tags auth
// @Produce json
// @Success 200 {object} dto.MessageResponse
// @Failure 401 {object} dto.ErrorResponse
// @Security CookieAuth
// @Security CSRFToken
// @Router /users/logout [post]
func (a *UserAPI) Logout(c *gin.Context) {
	// Logout 通过 service 层触发副作用
	//a.service.Logout() 暂时无逻辑
	a.sm.DestroySession(c.Writer)
	respondMessage(c, http.StatusOK, "logged out")
}
