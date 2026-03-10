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
// Author ownership is assigned by the service layer from the authenticated actor,
// not trusted from the request binding layer.
func (r *CreatePostRequest) ToEntity() *entity.Post {
	post := &entity.Post{
		Title:      r.Title,
		Content:    r.Content,
		Cover:      r.Cover,
		CategoryID: r.CategoryID,
		Status:     entity.StatusDraft, // 默认创建为草稿
	}
	if r.Tags != nil {
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
	// Status 由专用发布工作流接口管理：
	// POST /admin/posts/:id/publish 与 POST /admin/posts/:id/draft。
	// 这里保留字段兼容旧调用方，但 ToEntity 会显式忽略它。
	Status *int `json:"status"`
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
	// 状态切换必须走专用后台工作流接口，避免普通更新绕过业务约束。
	return post
}

// ToPatch converts the update DTO into a domain patch object.
// This keeps HTTP pointer semantics out of the service signature while preserving
// the distinction between omitted and explicitly provided scalar fields.
func (r *UpdatePostRequest) ToPatch() entity.PostPatch {
	patch := entity.PostPatch{
		Title:      r.Title,
		Content:    r.Content,
		Cover:      r.Cover,
		CategoryID: r.CategoryID,
	}
	if r.Tags != nil {
		patch.Tags = tagsFromIDs(r.Tags)
	}
	return patch
}

// tagsFromIDs converts a slice of tag IDs to a slice of entity.Tag.
// It centralizes the mapping logic to avoid duplication between DTO methods.
func tagsFromIDs(ids []uint) []entity.Tag {
	if ids == nil {
		return nil
	}
	tags := make([]entity.Tag, 0, len(ids))
	for _, id := range ids {
		tags = append(tags, entity.Tag{ID: id})
	}
	return tags
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
