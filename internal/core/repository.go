package core

import "KaldalisCMS/internal/core/entity"

type PostRepository interface {
	GetByID(id int) (entity.Post, error)
	Create(post entity.Post) (error)
	Update(post entity.Post) (error)
	Delete(id int) error
	GetAll() ([]entity.Post, error)
}
