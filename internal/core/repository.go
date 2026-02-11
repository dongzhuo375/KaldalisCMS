package core

import (
	"KaldalisCMS/internal/core/entity"
	"context"
)

type PostRepository interface {
	GetByID(id uint) (entity.Post, error)
	Create(post entity.Post) (entity.Post, error)
	Update(post entity.Post) error
	Delete(id uint) error
	GetAll() ([]entity.Post, error)
	IsSlugExists(slug string) (bool, error)
}

// MediaRepository defines persistence operations for media assets and post-media relations.
// Service layer should depend on this interface, not a specific DB implementation.
type MediaRepository interface {
	Create(ctx context.Context, asset *entity.MediaAsset) error
	GetByID(ctx context.Context, id uint) (entity.MediaAsset, error)
	List(ctx context.Context, ownerUserID *uint, offset, limit int, q string) ([]entity.MediaAsset, int64, error)
	Delete(ctx context.Context, id uint) error
	CountReferences(ctx context.Context, assetID uint) (int64, error)
	UpsertPostReferences(ctx context.Context, postID uint, purpose string, assetIDs []uint) error
	ListPostMedia(ctx context.Context, postID uint, purpose *string) ([]entity.MediaAsset, error)
	UpdateAssetFields(ctx context.Context, assetID uint, fields map[string]any) error
}

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	GetAll(ctx context.Context) ([]entity.User, error)
	GetByID(ctx context.Context, id uint) (entity.User, error)
	GetByUsername(ctx context.Context, username string) (entity.User, error)
	Create(ctx context.Context, user entity.User) error
	Update(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id uint) error
}

// TagRepository defines the interface for tag persistence.
type TagRepository interface {
	Create(ctx context.Context, tag entity.Tag) (entity.Tag, error)
	GetAll(ctx context.Context) ([]entity.Tag, error)
	GetByID(ctx context.Context, id uint) (entity.Tag, error)
	GetByName(ctx context.Context, name string) (entity.Tag, error)
	Update(ctx context.Context, tag entity.Tag) (entity.Tag, error)
	Delete(ctx context.Context, id uint) error
}
