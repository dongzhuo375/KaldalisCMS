package service

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	"context"
)

type UserService struct {
	repo core.UserRepository
	//enforcer *casbin.CachedEnforcer // 注入 Casbin 执行器
	//rdb      *redis.Client          // 注入共享 Redis
}

func NewUserService(repo core.UserRepository) *UserService {
	return &UserService{
		repo: repo,
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

func (s *UserService) Login(ctx context.Context, username, password string) (entity.User, error) {
	// 账号密码核对
	user, err := s.VerifyUser(username, password)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

// Logout 登出逻辑
func (s *UserService) Logout() {
	//逻辑留空
}
