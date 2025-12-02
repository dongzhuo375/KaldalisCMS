package v1

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAPI struct {
	userService core.UserService
}

func NewUserAPI(userService core.UserService) *UserAPI {
	return &UserAPI{
		userService: userService,
	}
}

// RegisterRoutes registers the user-related routes to the Gin router.
func (a *UserAPI) RegisterRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/register", a.Register)
		userRoutes.POST("/login", a.Login)
	}
}

// RegisterRequest defines the expected JSON structure for a registration request.
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Register handles new user registration.
func (a *UserAPI) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user", // Default role
	}

	if err := a.userService.CreateUser(user); err != nil {
		if errors.Is(err, core.ErrDuplicate) {
			c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// LoginRequest defines the expected JSON structure for a login request.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles user authentication.
func (a *UserAPI) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.userService.VerifyUser(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) || errors.Is(err, core.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		return
	}

	// In a real application, you would generate a JWT token here.
	// For now, we'll just return a success message and the user's data (excluding password).
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}