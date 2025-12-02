package core

import "KaldalisCMS/internal/core/entity"

type PostRepository interface {
	GetByID(id int) (entity.Post, error)
	Create(post entity.Post) error
	Update(post entity.Post) error
	Delete(id int) error
	GetAll() ([]entity.Post, error)
	IsSlugExists(slug string) (bool, error)
}

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	GetAll() ([]entity.User, error)
	GetByID(id int) (entity.User, error)
	GetByUsername(username string) (entity.User, error)
	Create(user entity.User) error
	Update(user entity.User) error
	Delete(id int) error
}
