package dto

// UserRegisterRequest defines the request body for user registration.
// It only includes fields that should be provided by the user.
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}
