package dto

// import (
// 	"KaldalisCMS/internal/core/entity" // Import entity package
// 	"time"                              // Import time package
// )

// UserRegisterRequest defines the request body for user registration.
// It only includes fields that should be provided by the user.
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

// // UserLoginRequest defines the request body for user login.
// type UserLoginRequest struct {
// 	Email    string `json:"email" binding:"required,email"`
// 	Password string `json:"password" binding:"required"`
// }

// // UserResponse defines the data structure for sending user information back to the client.
// type UserResponse struct {
// 	ID        uint      `json:"id"`
// 	Username  string    `json:"username"`
// 	Email     string    `json:"email"`
// 	Role      string    `json:"role"`
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`
// }

// // UpdateUserRequest defines the request body for updating user information.
// // Fields are pointers to allow partial updates (PATCH).
// type UpdateUserRequest struct {
// 	Username *string `json:"username" binding:"omitempty,min=3,max=50"`
// 	Email    *string `json:"email" binding:"omitempty,email"`
// 	Password *string `json:"password" binding:"omitempty,min=6"`
// 	Role     *string `json:"role" binding:"omitempty"`
// }

// // ToUserResponse converts an entity.User to a UserResponse DTO.
// func ToUserResponse(user *entity.User) *UserResponse {
// 	if user == nil {
// 		return nil
// 	}
// 	return &UserResponse{
// 		ID:        user.ID,
// 		Username:  user.Username,
// 		Email:     user.Email,
// 		Role:      user.Role,
// 		CreatedAt: user.CreatedAt,
// 		UpdatedAt: user.UpdatedAt,
// 	}
// }
