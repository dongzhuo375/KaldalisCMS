package service

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity" // Service 只能使用 Entity
)

type UserService struct {
	repo core.UserRepository
}

func NewUserService(repo core.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}
