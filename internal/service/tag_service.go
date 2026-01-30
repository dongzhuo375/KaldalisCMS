package service

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	"context"
	"errors"
	"strings"
)

// tagService implements core.TagService.
type tagService struct {
	repo core.TagRepository
}

// NewTagService creates a TagService.
func NewTagService(repo core.TagRepository) core.TagService {
	return &tagService{repo: repo}
}

// Create creates a new tag.
func (s *tagService) Create(ctx context.Context, tag entity.Tag) (entity.Tag, error) {
	tag.Name = strings.TrimSpace(tag.Name)
	if tag.Name == "" {
		return entity.Tag{}, core.ErrInvalidInput
	}

	// Best-effort uniqueness check by name.
	// Expected repo contract:
	//   - return (Tag, nil) when found
	//   - return (zero, core.ErrNotFound) when not found
	existing, err := s.repo.GetByName(ctx, tag.Name)
	switch {
	case err == nil:
		if existing.ID != 0 {
			return entity.Tag{}, core.ErrDuplicate
		}
	case errors.Is(err, core.ErrNotFound):
		// ok, not exists
	default:
		// unexpected repo error
		return entity.Tag{}, err
	}

	return s.repo.Create(ctx, tag)
}

// GetAll returns all tags.
func (s *tagService) GetAll(ctx context.Context) ([]entity.Tag, error) {
	// No business-specific validation here; just delegate to repository.
	return s.repo.GetAll(ctx)
}

// GetByID returns a tag by ID.
func (s *tagService) GetByID(ctx context.Context, id uint) (entity.Tag, error) {
	if id == 0 {
		return entity.Tag{}, core.ErrInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}

// GetByName returns a tag by name.
func (s *tagService) GetByName(ctx context.Context, name string) (entity.Tag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return entity.Tag{}, core.ErrInvalidInput
	}
	return s.repo.GetByName(ctx, name)
}

// Update updates a tag.
func (s *tagService) Update(ctx context.Context, tag entity.Tag) (entity.Tag, error) {
	if tag.ID == 0 {
		return entity.Tag{}, core.ErrInvalidInput
	}
	if tag.Name != "" {
		tag.Name = strings.TrimSpace(tag.Name)
		if tag.Name == "" {
			return entity.Tag{}, core.ErrInvalidInput
		}
	}

	return s.repo.Update(ctx, tag)
}

// Delete deletes a tag by ID.
func (s *tagService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return core.ErrInvalidInput
	}
	return s.repo.Delete(ctx, id)
}
