package core

import (
	"KaldalisCMS/internal/core/entity"
	"context"
)

type PostRepository interface {
	GetByID(id uint) (entity.Post, error)
	Create(post entity.Post) error
	Update(post entity.Post) error
	Delete(id uint) error
	GetAll() ([]entity.Post, error)
	IsSlugExists(slug string) (bool, error)
}

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	GetAll(ctx context.Context) ([]entity.User, error)
	GetByID(ctx context.Context, id uint) (entity.User, error)
	GetByUsername(ctx context.Context, username string) (entity.User, error)
	Create(ctx context.Context, user entity.User) error
	Update(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id uint) error
}

// TagRepository defines the interface for tag persistence.
type TagRepository interface {
	Create(ctx context.Context, tag entity.Tag) (entity.Tag, error)
	GetAll(ctx context.Context) ([]entity.Tag, error)
	GetByID(ctx context.Context, id uint) (entity.Tag, error)
	GetByName(ctx context.Context, name string) (entity.Tag, error)
	Update(ctx context.Context, tag entity.Tag) (entity.Tag, error)
	Delete(ctx context.Context, id uint) error
}
