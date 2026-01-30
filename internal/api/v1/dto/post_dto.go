package dto

import (
	"KaldalisCMS/internal/core/entity"
	"time"
)

// CreatePostRequest defines the structure for creating a new post.
type CreatePostRequest struct {
	Title      string `json:"title" binding:"required,min=1,max=100"`
	Content    string `json:"content"`
	Cover      string `json:"cover" binding:"max=255"`
	CategoryID *uint  `json:"category_id"`
	Tags       []uint `json:"tags"`
}

// ToEntity converts a CreatePostRequest DTO to an entity.Post.
// It requires the authorID to be passed as it's not part of the request body.
func (r *CreatePostRequest) ToEntity(authorID uint) *entity.Post {
	post := &entity.Post{
		Title:      r.Title,
		Content:    r.Content,
		Cover:      r.Cover,
		CategoryID: r.CategoryID,
		AuthorID:   authorID,           // <-- Add AuthorID
		Status:     entity.StatusDraft, // 默认创建为草稿
	}
	if r.Tags != nil {
		// Post.Tags 的元素类型在 entity 包内是 Tag（同包类型），不需要用 entity.Tag 显式限定。
		post.Tags = make([]entity.Tag, 0, len(r.Tags))
		for _, tagID := range r.Tags {
			post.Tags = append(post.Tags, entity.Tag{ID: tagID})
		}
	}
	return post
}

// UpdatePostRequest defines the structure for updating an existing post.
type UpdatePostRequest struct {
	Title      *string `json:"title" binding:"omitempty,min=1,max=100"`
	Content    *string `json:"content"`
	Cover      *string `json:"cover" binding:"omitempty,max=255"`
	CategoryID *uint   `json:"category_id"`
	Tags       []uint  `json:"tags"`
	Status     *int    `json:"status"`
}

// ToEntity creates and returns a new entity.Post from the UpdatePostRequest.
// Only non-nil fields in the DTO will be set in the new entity.
func (r *UpdatePostRequest) ToEntity() entity.Post {
	post := entity.Post{} // Create a new entity

	if r.Title != nil {
		post.Title = *r.Title
	}
	if r.Content != nil {
		post.Content = *r.Content
	}
	if r.Cover != nil {
		post.Cover = *r.Cover
	}
	if r.CategoryID != nil {
		post.CategoryID = r.CategoryID
	}
	if r.Tags != nil {
		post.Tags = make([]entity.Tag, 0, len(r.Tags))
		for _, tagID := range r.Tags {
			post.Tags = append(post.Tags, entity.Tag{ID: tagID})
		}
	}
	if r.Status != nil {
		post.Status = *r.Status
	}
	return post
}

// --- 以下是建议新增和修改的部分 ---

// PostResponse is the DTO for a single post.
type PostResponse struct {
	ID        uint              `json:"id"`
	Title     string            `json:"title"`
	Slug      string            `json:"slug"`
	Content   string            `json:"content"`
	Cover     string            `json:"cover"`
	Status    int               `json:"status"`
	Author    AuthorResponse    `json:"author"`
	Category  *CategoryResponse `json:"category,omitempty"`
	Tags      []TagResponse     `json:"tags,omitempty"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

// AuthorResponse is the DTO for post author.
type AuthorResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// CategoryResponse is the DTO for post category.
type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ToPostResponse converts an entity.Post to a PostResponse DTO.
func ToPostResponse(post *entity.Post) *PostResponse {
	if post == nil {
		return nil
	}

	res := &PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Slug:      post.Slug,
		Content:   post.Content,
		Cover:     post.Cover,
		Status:    post.Status,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
		Author: AuthorResponse{
			ID:       post.Author.ID,
			Username: post.Author.Username,
		},
	}

	if post.CategoryID != nil {
		res.Category = &CategoryResponse{
			ID:   post.Category.ID,
			Name: post.Category.Name,
		}
	}

	if len(post.Tags) > 0 {
		res.Tags = make([]TagResponse, len(post.Tags))
		for i, tag := range post.Tags {
			res.Tags[i] = TagResponse{
				ID:   tag.ID,
				Name: tag.Name,
			}
		}
	}

	return res
}

// ToPostListResponse converts a slice of entity.Post to a slice of PostResponse DTOs.
func ToPostListResponse(posts []entity.Post) []*PostResponse {
	if len(posts) == 0 {
		return []*PostResponse{}
	}

	res := make([]*PostResponse, len(posts))
	for i, post := range posts {
		res[i] = ToPostResponse(&post)
	}
	return res
}
