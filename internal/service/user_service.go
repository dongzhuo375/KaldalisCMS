package service

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	"KaldalisCMS/internal/infra/auth"
	"net/http"
)

type UserService struct {
	repo    core.UserRepository
	authCfg auth.Config
	//enforcer *casbin.CachedEnforcer // 注入 Casbin 执行器
	//rdb      *redis.Client          // 注入共享 Redis
}

func NewUserService(repo core.UserRepository, cfg auth.Config) *UserService {
	return &UserService{
		repo:    repo,
		authCfg: cfg,
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

func (s *UserService) Login(w http.ResponseWriter, username, password string) (entity.User, error) {
	user, err := s.VerifyUser(username, password)
	if err != nil {
		return entity.User{}, err
	}

	// 调用 infra 包的公开函数，并传入 s.authCfg
	if err := auth.EstablishSession(w, s.authCfg, user.ID); err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// Logout 登出逻辑
func (s *UserService) Logout(w http.ResponseWriter) {
	auth.DestroySession(w, s.authCfg)
}
