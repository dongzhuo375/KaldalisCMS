package core

import "KaldalisCMS/internal/model"

type PostService interface {
	GetAllPosts() ([]model.Post, error)
	GetPostByID(id int) (model.Post, error)
	CreatePost(post model.Post) (model.Post, error)
	UpdatePost(id int, post model.Post) (model.Post, error)
	DeletePost(id int) error
}
