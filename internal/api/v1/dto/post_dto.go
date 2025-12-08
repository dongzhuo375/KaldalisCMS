package dto

import (
	"KaldalisCMS/internal/core/entity"
)

// CreatePostRequest defines the structure for creating a new post.
type CreatePostRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content"`
	Cover      string `json:"cover"`
	CategoryID *uint  `json:"category_id"`
	Tags       []uint `json:"tags"`
}

// ToEntity converts a CreatePostRequest DTO to an entity.Post.
func (r *CreatePostRequest) ToEntity() *entity.Post {
	post := &entity.Post{
		Title:      r.Title,
		Content:    r.Content,
		Cover:      r.Cover,
		CategoryID: r.CategoryID,
	}
	if r.Tags != nil {
		post.Tags = make([]entity.Tag, len(r.Tags))
		for i, tagID := range r.Tags {
			post.Tags[i] = entity.Tag{ID: tagID}
		}
	}
	return post
}

// UpdatePostRequest defines the structure for updating an existing post.
type UpdatePostRequest struct {
	Title      *string `json:"title"`
	Content    *string `json:"content"`
	Cover      *string `json:"cover"`
	CategoryID *uint   `json:"category_id"`
	Tags       []uint  `json:"tags"`
	Status     *int    `json:"status"`
}

