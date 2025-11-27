package service

import (
	"KaldalisCMS/internal/model"
	"KaldalisCMS/internal/repository"
)

type PostService struct {
	repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) GetAllPosts() ([]model.Post, error) {
	return s.repo.GetAll()
}

func (s *PostService) GetPostByID(id int) (model.Post, error) {
	return s.repo.GetByID(id)
}

func (s *PostService) CreatePost(post model.Post) (model.Post, error) {
	// In a real application, you might add validation or other business logic here.
	return s.repo.Create(post)
}

func (s *PostService) UpdatePost(id int, post model.Post) (model.Post, error) {
	return s.repo.Update(id, post)
}

func (s *PostService) DeletePost(id int) error {
	return s.repo.Delete(id)
}
