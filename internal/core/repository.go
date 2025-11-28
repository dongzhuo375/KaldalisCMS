package core

import "KaldalisCMS/internal/model"

type PostRepository interface {
	Get(id int) (*model.Post, error)
	Create(post *model.Post) error
	Update(post *model.Post) error
	Delete(id int) error
	List() ([]*model.Post, error)
}
