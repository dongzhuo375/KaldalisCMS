package dto

import "KaldalisCMS/internal/core/entity"

// CreateTagRequest defines the request body for creating a tag.
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=50"`
}

func (r *CreateTagRequest) ToEntity() entity.Tag {
	return entity.Tag{Name: r.Name}
}

// UpdateTagRequest defines the request body for updating a tag.
type UpdateTagRequest struct {
	Name *string `json:"name" binding:"omitempty,min=1,max=50"`
}

func (r *UpdateTagRequest) ToEntity() entity.Tag {
	t := entity.Tag{}
	if r.Name != nil {
		t.Name = *r.Name
	}
	return t
}

// TagResponse is the DTO for returning tag info.
type TagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func ToTagResponse(tag entity.Tag) TagResponse {
	return TagResponse{ID: tag.ID, Name: tag.Name}
}

func ToTagListResponse(tags []entity.Tag) []TagResponse {
	if len(tags) == 0 {
		return []TagResponse{}
	}
	res := make([]TagResponse, len(tags))
	for i, t := range tags {
		res[i] = ToTagResponse(t)
	}
	return res
}
