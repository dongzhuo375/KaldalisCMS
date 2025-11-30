package core

import "KaldalisCMS/internal/core/entity"

type PostRepository interface {
	GetByID(id int) (entity.Post, error)
	Create(post entity.Post) (entity.Post, error)
	Update(post entity.Post) (entity.Post, error)
	Delete(id int) error
	GetAll() ([]entity.Post, error)
	FindByID(id int) (entity.Post, error)
}
