package repository

import (
	"KaldalisCMS/internal/model"
	"fmt"
)

type PostRepository struct{}

func NewPostRepository() *PostRepository {
	return &PostRepository{}
}

func (r *PostRepository) GetAll() ([]model.Post, error) {
	InMemoryDB.RLock()
	defer InMemoryDB.RUnlock()

	posts := make([]model.Post, 0, len(InMemoryDB.posts))
	for _, post := range InMemoryDB.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) GetByID(id int) (model.Post, error) {
	InMemoryDB.RLock()
	defer InMemoryDB.RUnlock()

	post, exists := InMemoryDB.posts[id]
	if !exists {
		return model.Post{}, fmt.Errorf("post with ID %d not found", id)
	}
	return post, nil
}

func (r *PostRepository) Create(post model.Post) (model.Post, error) {
	InMemoryDB.Lock()
	defer InMemoryDB.Unlock()

	InMemoryDB.counter++
	post.ID = InMemoryDB.counter
	InMemoryDB.posts[post.ID] = post
	return post, nil
}

func (r *PostRepository) Update(id int, post model.Post) (model.Post, error) {
	InMemoryDB.Lock()
	defer InMemoryDB.Unlock()

	_, exists := InMemoryDB.posts[id]
	if !exists {
		return model.Post{}, fmt.Errorf("post with ID %d not found", id)
	}
	post.ID = id
	InMemoryDB.posts[id] = post
	return post, nil
}

func (r *PostRepository) Delete(id int) error {
	InMemoryDB.Lock()
	defer InMemoryDB.Unlock()

	_, exists := InMemoryDB.posts[id]
	if !exists {
		return fmt.Errorf("post with ID %d not found", id)
	}
	delete(InMemoryDB.posts, id)
	return nil
}
