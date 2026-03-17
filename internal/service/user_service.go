package service

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	"context"
	"errors"
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
func (s *UserService) CreateUser(ctx context.Context, user entity.User) error {
	// The user.Password field currently holds the plaintext password.
	// We use the entity's own method to hash it.
	if err := user.SetPassword(user.Password); err != nil {
		return core.ErrInvalidInput
	}

	// Now user.Password holds the hashed password.
	// We can pass the entity to the repository to be created.
	return normalizeServiceErrorWithOpMsg("user.create", "create user failed", s.repo.Create(ctx, user))
}

// VerifyUser 只做验证并返回 user
func (s *UserService) VerifyUser(ctx context.Context, username, password string) (entity.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return entity.User{}, core.ErrInvalidCredentials
		}
		return entity.User{}, normalizeServiceErrorWithOpMsg("user.verify.lookup", "lookup user by username failed", err)
	}
	if !user.CheckPassword(password) {
		return entity.User{}, core.ErrInvalidCredentials
	}
	return user, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (entity.User, error) {
	// 账号密码核对
	user, err := s.VerifyUser(ctx, username, password)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

// Logout 登出逻辑
func (s *UserService) Logout() {
	//逻辑留空
}
