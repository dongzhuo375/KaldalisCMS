package repository

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	"KaldalisCMS/internal/infra/model"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Ensure *MediaRepository implements core.MediaRepository.
var _ core.MediaRepository = (*MediaRepository)(nil)

// --- Mapper helpers ---

func mediaModelToEntity(m model.MediaAsset) entity.MediaAsset {
	return entity.MediaAsset{
		ID:           m.ID,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		OwnerUserID:  m.OwnerUserID,
		OriginalName: m.OriginalName,
		StoredName:   m.StoredName,
		Ext:          m.Ext,
		MimeType:     m.MimeType,
		SizeBytes:    m.SizeBytes,
		SHA256:       m.SHA256,
		Storage:      m.Storage,
		ObjectKey:    m.ObjectKey,
		Url:          m.Url,
		Width:        m.Width,
		Height:       m.Height,
	}
}

func mediaEntityToModel(e entity.MediaAsset) model.MediaAsset {
	return model.MediaAsset{
		ID:           e.ID,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
		OwnerUserID:  e.OwnerUserID,
		OriginalName: e.OriginalName,
		StoredName:   e.StoredName,
		Ext:          e.Ext,
		MimeType:     e.MimeType,
		SizeBytes:    e.SizeBytes,
		SHA256:       e.SHA256,
		Storage:      e.Storage,
		ObjectKey:    e.ObjectKey,
		Url:          e.Url,
		Width:        e.Width,
		Height:       e.Height,
	}
}

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

var ErrMediaNotFound = errors.New("media asset not found")

func (r *MediaRepository) Create(ctx context.Context, asset *entity.MediaAsset) error {
	if asset == nil {
		return fmt.Errorf("media_repository.Create: asset is nil")
	}
	m := mediaEntityToModel(*asset)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("media_repository.Create: %w", err)
	}
	asset.ID = m.ID
	asset.CreatedAt = m.CreatedAt
	asset.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *MediaRepository) GetByID(ctx context.Context, id uint) (entity.MediaAsset, error) {
	var m model.MediaAsset
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.MediaAsset{}, ErrMediaNotFound
		}
		return entity.MediaAsset{}, fmt.Errorf("media_repository.GetByID: %w", err)
	}
	return mediaModelToEntity(m), nil
}

func (r *MediaRepository) List(ctx context.Context, ownerUserID *uint, offset, limit int, q string) ([]entity.MediaAsset, int64, error) {
	var ms []model.MediaAsset
	query := r.db.WithContext(ctx).Model(&model.MediaAsset{})
	if ownerUserID != nil {
		query = query.Where("owner_user_id = ?", *ownerUserID)
	}
	if q != "" {
		query = query.Where("original_name ILIKE ?", "%"+q+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("media_repository.List.count: %w", err)
	}
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&ms).Error; err != nil {
		return nil, 0, fmt.Errorf("media_repository.List: %w", err)
	}

	out := make([]entity.MediaAsset, 0, len(ms))
	for _, m := range ms {
		out = append(out, mediaModelToEntity(m))
	}
	return out, total, nil
}

func (r *MediaRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.MediaAsset{}, id).Error; err != nil {
		return fmt.Errorf("media_repository.Delete: %w", err)
	}
	return nil
}

func (r *MediaRepository) CountReferences(ctx context.Context, assetID uint) (int64, error) {
	var cnt int64
	if err := r.db.WithContext(ctx).Model(&model.PostAsset{}).Where("asset_id = ?", assetID).Count(&cnt).Error; err != nil {
		return 0, fmt.Errorf("media_repository.CountReferences: %w", err)
	}
	return cnt, nil
}

func (r *MediaRepository) UpsertPostReferences(ctx context.Context, postID uint, purpose string, assetIDs []uint) error {
	if purpose == "" {
		purpose = "content"
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("post_id = ? AND purpose = ?", postID, purpose).Delete(&model.PostAsset{}).Error; err != nil {
			return fmt.Errorf("media_repository.UpsertPostReferences.delete: %w", err)
		}
		if len(assetIDs) == 0 {
			return nil
		}
		rows := make([]model.PostAsset, 0, len(assetIDs))
		for _, id := range assetIDs {
			rows = append(rows, model.PostAsset{PostID: postID, AssetID: id, Purpose: purpose})
		}
		if err := tx.Create(&rows).Error; err != nil {
			return fmt.Errorf("media_repository.UpsertPostReferences.insert: %w", err)
		}
		return nil
	})
}

func (r *MediaRepository) ListPostMedia(ctx context.Context, postID uint, purpose *string) ([]entity.MediaAsset, error) {
	q := r.db.WithContext(ctx).
		Table("media_assets").
		Select("media_assets.*").
		Joins("JOIN post_assets ON post_assets.asset_id = media_assets.id").
		Where("post_assets.post_id = ?", postID)
	if purpose != nil {
		q = q.Where("post_assets.purpose = ?", *purpose)
	}
	q = q.Order("media_assets.created_at DESC")

	var ms []model.MediaAsset
	if err := q.Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("media_repository.ListPostMedia: %w", err)
	}
	out := make([]entity.MediaAsset, 0, len(ms))
	for _, m := range ms {
		out = append(out, mediaModelToEntity(m))
	}
	return out, nil
}

func (r *MediaRepository) UpdateAssetFields(ctx context.Context, assetID uint, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).Model(&model.MediaAsset{}).Where("id = ?", assetID).Updates(fields).Error; err != nil {
		return fmt.Errorf("media_repository.UpdateAssetFields: %w", err)
	}
	return nil
}

// Db returns the underlying gorm DB handle.
// Intended for small internal operations that don't warrant dedicated repository methods.
// Deprecated: prefer adding explicit repository methods (e.g., UpdateAssetFields).
func (r *MediaRepository) Db() *gorm.DB {
	return r.db
}
