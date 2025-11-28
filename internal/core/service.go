package core

import "KaldalisCMS/internal/model"

type PostService interface {
	GetPost(id int) (*model.Post, error)
	CreatePost(post *model.Post) error
	UpdatePost(post *model.Post) error
	DeletePost(id int) error
	ListPosts() ([]*model.Post, error)
}
