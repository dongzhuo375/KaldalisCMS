package service

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	"net/http"
)

type UserService struct {
	repo    core.UserRepository
	authMgr core.AuthManager
}

func NewUserService(repo core.UserRepository, authMgr core.AuthManager) *UserService {
	return &UserService{
		repo:    repo,
		authMgr: authMgr,
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

// VerifyUser 只做验证并返回 user
func (s *UserService) VerifyUser(username, password string) (entity.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return entity.User{}, err
	}
	if !user.CheckPassword(password) {
		return entity.User{}, core.ErrInvalidCredentials
	}
	return user, nil
}

// Login 封装验证 + infra auth 登录副作用（写 cookie）
func (s *UserService) Login(w http.ResponseWriter, username, password string, secureFlag bool) (entity.User, error) {
	user, err := s.VerifyUser(username, password)
	if err != nil {
		return entity.User{}, err
	}
	// 使用注入的 AuthManager 写 cookie
	if err := s.authMgr.Login(w, user.ID, secureFlag); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

// Logout 登出作用
func (s *UserService) Logout(w http.ResponseWriter) {
	s.authMgr.Logout(w)
}
