package core

import (
	"KaldalisCMS/internal/core/entity"
	"context"
)

type PostService interface {
	GetAllPosts() ([]entity.Post, error)
	GetPostByID(id uint) (entity.Post, error)
	CreatePost(post entity.Post) error
	UpdatePost(id uint, post entity.Post) error
	DeletePost(id uint) error
	PublishPost(id uint) error
	DraftPost(id uint) error
}

type UserService interface {
	CreateUser(ctx context.Context, user entity.User) error
	VerifyUser(ctx context.Context, username, password string) (entity.User, error)
	Login(ctx context.Context, username, password string) (entity.User, error)
	Logout()
	//后面估计还加
}
