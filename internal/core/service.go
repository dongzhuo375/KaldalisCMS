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

// TagService defines tag-related business operations.
type TagService interface {
	Create(ctx context.Context, tag entity.Tag) (entity.Tag, error)
	GetAll(ctx context.Context) ([]entity.Tag, error)
	GetByID(ctx context.Context, id uint) (entity.Tag, error)
	GetByName(ctx context.Context, name string) (entity.Tag, error)
	Update(ctx context.Context, tag entity.Tag) (entity.Tag, error)
	Delete(ctx context.Context, id uint) error
}
