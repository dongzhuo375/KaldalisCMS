package repository

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	"context"
	"strings"
	"sync"
)

// InMemoryTagRepository is a temporary repository implementation for tests/smoke runs.
//
// NOTE: This is NOT intended for production use.
type InMemoryTagRepository struct {
	mu     sync.RWMutex
	nextID uint
	byID   map[uint]entity.Tag
}

var _ core.TagRepository = (*InMemoryTagRepository)(nil)

func NewInMemoryTagRepository() *InMemoryTagRepository {
	return &InMemoryTagRepository{
		nextID: 1,
		byID:   make(map[uint]entity.Tag),
	}
}

func (r *InMemoryTagRepository) Create(ctx context.Context, tag entity.Tag) (entity.Tag, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	tag.Name = strings.TrimSpace(tag.Name)
	if tag.Name == "" {
		return entity.Tag{}, core.ErrInvalidInput
	}

	// uniqueness check by name
	for _, existing := range r.byID {
		if strings.EqualFold(existing.Name, tag.Name) {
			return entity.Tag{}, core.ErrDuplicate
		}
	}

	tag.ID = r.nextID
	r.nextID++
	r.byID[tag.ID] = tag
	return tag, nil
}

func (r *InMemoryTagRepository) GetAll(ctx context.Context) ([]entity.Tag, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]entity.Tag, 0, len(r.byID))
	for _, t := range r.byID {
		res = append(res, t)
	}
	return res, nil
}

func (r *InMemoryTagRepository) GetByID(ctx context.Context, id uint) (entity.Tag, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.byID[id]
	if !ok {
		return entity.Tag{}, core.ErrNotFound
	}
	return t, nil
}

func (r *InMemoryTagRepository) GetByName(ctx context.Context, name string) (entity.Tag, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	name = strings.TrimSpace(name)
	if name == "" {
		return entity.Tag{}, core.ErrInvalidInput
	}
	for _, t := range r.byID {
		if strings.EqualFold(t.Name, name) {
			return t, nil
		}
	}
	return entity.Tag{}, core.ErrNotFound
}

func (r *InMemoryTagRepository) Update(ctx context.Context, tag entity.Tag) (entity.Tag, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	if tag.ID == 0 {
		return entity.Tag{}, core.ErrInvalidInput
	}

	existing, ok := r.byID[tag.ID]
	if !ok {
		return entity.Tag{}, core.ErrNotFound
	}

	// only update Name when provided
	if tag.Name != "" {
		name := strings.TrimSpace(tag.Name)
		if name == "" {
			return entity.Tag{}, core.ErrInvalidInput
		}

		for id, t := range r.byID {
			if id != tag.ID && strings.EqualFold(t.Name, name) {
				return entity.Tag{}, core.ErrDuplicate
			}
		}

		existing.Name = name
	}

	r.byID[tag.ID] = existing
	return existing, nil
}

func (r *InMemoryTagRepository) Delete(ctx context.Context, id uint) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[id]; !ok {
		return core.ErrNotFound
	}
	delete(r.byID, id)
	return nil
}
