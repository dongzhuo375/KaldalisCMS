package service

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserService struct {
	repo               core.UserRepository
	jwtSecret          string
	jwtExpirationHours int
}

func NewUserService(repo core.UserRepository, jwtSecret string, jwtExpirationHours int) *UserService {
	return &UserService{
		repo:               repo,
		jwtSecret:          jwtSecret,
		jwtExpirationHours: jwtExpirationHours,
	}
}

// CreateUser handles the business logic for creating a new user.
// It hashes the password before passing the data to the repository.
func (s *UserService) CreateUser(user entity.User) error {
	// The user.Password field currently holds the plaintext password.
	// We use the entity's own method to hash it.
	if err := user.SetPassword(user.Password); err != nil {
		return err
	}

	// Now user.Password holds the hashed password.
	// We can pass the entity to the repository to be created.
	return s.repo.Create(user)
}

// VerifyUser checks a user's credentials.
// It retrieves the user by username and then validates the password.
func (s *UserService) VerifyUser(username, password string) (entity.User, string, error) {
	// Retrieve the user from the database.
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		// This could be ErrNotFound or a database connection error.
		// The API layer will handle them accordingly.
		return entity.User{}, "", err
	}

	// Use the entity's method to check the password.
	if !user.CheckPassword(password) {
		return entity.User{}, "", core.ErrInvalidCredentials
	}

	// Password is correct. Generate a JWT token.
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * time.Duration(s.jwtExpirationHours)).Unix(),
		"iat":    time.Now().Unix(),
		"iss":    "KaldalisCMS",
	})

	// Sign the token with the secret.
	tokenString, err := claims.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return entity.User{}, "", err
	}

	return user, tokenString, nil
}
