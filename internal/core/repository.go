package core

import (
	"KaldalisCMS/internal/core/entity"
	"context"
	"time"
)

type PostRepository interface {
	GetByID(ctx context.Context, id uint) (entity.Post, error)
	GetPublishedByID(ctx context.Context, id uint) (entity.Post, error)
	GetDraftByIDAndAuthor(ctx context.Context, id uint, authorID uint) (entity.Post, error)
	Create(ctx context.Context, post entity.Post) (entity.Post, error)
	Update(ctx context.Context, post entity.Post) error
	Delete(ctx context.Context, id uint) error
	GetAll(ctx context.Context) ([]entity.Post, error)
	GetPublished(ctx context.Context) ([]entity.Post, error)
	GetDraftsByAuthor(ctx context.Context, authorID uint) ([]entity.Post, error)
	IsSlugExists(ctx context.Context, slug string) (bool, error)
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
	UpdateStatus(ctx context.Context, id uint, status entity.MediaStatus) error
	ListPendingOlderThan(ctx context.Context, cutoff time.Time, limit int) ([]entity.MediaAsset, error)
	ListSoftDeletedOlderThan(ctx context.Context, cutoff time.Time, limit int) ([]entity.MediaAsset, error)
	DeletePhysical(ctx context.Context, id uint) error
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
