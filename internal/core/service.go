package core

import "KaldalisCMS/internal/core/entity"

type PostService interface {
	GetAllPosts() ([]entity.Post, error)
	GetPostByID(id int) (entity.Post, error)
	CreatePost(post entity.Post) ( error)
	UpdatePost(id int, post entity.Post) ( error)
	DeletePost(id int) error
	PublishPost(id int) error
	DraftPost(id int) error
}

type UserService interface {
	CreateUser(user entity.User) error
	VerifyUser(username, password string) (entity.User, error)
	//后面估计还加
}


