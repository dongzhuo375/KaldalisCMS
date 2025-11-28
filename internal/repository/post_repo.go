package repository

import (
	"KaldalisCMS/internal/model"
)

type PostRepository struct{}

func NewPostRepository() *PostRepository {
	return &PostRepository{}
}

func (r *PostRepository) GetAll() ([]model.Post, error) {
	var posts []model.Post
	if err := DB.Preload("Author").Preload("Category").Preload("Tags").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) GetByID(id int) (model.Post, error) {
	var post model.Post
	if err := DB.Preload("Author").Preload("Category").Preload("Tags").First(&post, id).Error; err != nil {
		return model.Post{}, err
	}
	return post, nil
}

func (r *PostRepository) Create(post model.Post) (model.Post, error) {
	if err := DB.Create(&post).Error; err != nil {
		return model.Post{}, err
	}
	return post, nil
}

func (r *PostRepository) Update(id int, post model.Post) (model.Post, error) {
    var existingPost model.Post
    if err := DB.First(&existingPost, id).Error; err != nil {
        return model.Post{}, err
    }

	// It's generally better to update specific fields instead of the whole object
	// but for this refactoring, we'll update the whole object.
	// Note: GORM's Save updates all fields, which might not be what you want for partial updates.
	// Consider using .Model(&post).Updates(post) for more control.
	if err := DB.Model(&existingPost).Updates(post).Error; err != nil {
		return model.Post{}, err
	}
	return existingPost, nil
}

func (r *PostRepository) Delete(id int) error {
	if err := DB.Delete(&model.Post{}, id).Error; err != nil {
		return err
	}
	return nil
}
