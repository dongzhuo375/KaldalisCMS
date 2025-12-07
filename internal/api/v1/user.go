package v1

import (
	"KaldalisCMS/internal/core/entity"
	"KaldalisCMS/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAPI struct {
	service *service.UserService
}

func NewUserAPI(service *service.UserService) *UserAPI {
	return &UserAPI{
		service: service,
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
func (api *UserAPI) Register(c *gin.Context) {
	var newUser entity.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	
	if err := api.service.CreateUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}


// **dzcake,请您在下面实现jwt**
// Login handles user authentication.
func (a *UserAPI) Login(c *gin.Context) {
	// Use an anonymous struct for binding, as there is no specific DTO
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := a.service.VerifyUser(req.Username, req.Password)
	if err != nil {
		// Simplified error handling
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token, // Return the JWT token
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}