package dto

// UserLoginRequest defines login payload.
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginUserResponse is embedded in LoginResponse.
type LoginUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// LoginResponse defines login success payload.
type LoginResponse struct {
	Message   string            `json:"message"`
	User      LoginUserResponse `json:"user"`
	ExpiresAt string            `json:"expires_at"`
}
