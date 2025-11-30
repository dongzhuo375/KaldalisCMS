package core

import "KaldalisCMS/internal/model"

type PostRepository interface {
	GetByID(id int) (model.Post, error)
	Create(post model.Post) (model.Post, error)
	Update(id int, post model.Post) (model.Post, error)
	Delete(id int) error
	GetAll() ([]model.Post, error)
}
