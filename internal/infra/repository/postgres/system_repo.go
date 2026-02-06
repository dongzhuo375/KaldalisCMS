package repository

import (
	"KaldalisCMS/internal/infra/model"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SystemRepository provides DB access for system-level settings.
// It is intentionally small and focused: setup flow must be concurrency-safe.
type SystemRepository struct {
	db *gorm.DB
}

func NewSystemRepository(db *gorm.DB) *SystemRepository {
	return &SystemRepository{db: db}
}

var ErrSystemSettingNotFound = errors.New("system setting not found")

// EnsureSingletonRow guarantees that the singleton row (ID=1) exists.
// It is safe to call multiple times.
func (r *SystemRepository) EnsureSingletonRow(ctx context.Context, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}

	row := model.SystemSetting{ID: 1}
	// Insert only if missing.
	if err := tx.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&row).Error; err != nil {
		return fmt.Errorf("system_repository.EnsureSingletonRow: %w", err)
	}
	return nil
}

// Get returns the singleton settings row.
func (r *SystemRepository) Get(ctx context.Context) (model.SystemSetting, error) {
	var s model.SystemSetting
	err := r.db.WithContext(ctx).First(&s, 1).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.SystemSetting{}, ErrSystemSettingNotFound
		}
		return model.SystemSetting{}, fmt.Errorf("system_repository.Get: %w", err)
	}
	return s, nil
}

// LockForUpdate locks the singleton row (ID=1) within an existing transaction.
// Caller must pass a transactional *gorm.DB.
func (r *SystemRepository) LockForUpdate(ctx context.Context, tx *gorm.DB) (model.SystemSetting, error) {
	var s model.SystemSetting
	if tx == nil {
		return model.SystemSetting{}, fmt.Errorf("system_repository.LockForUpdate: tx is required")
	}

	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&s, 1).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.SystemSetting{}, ErrSystemSettingNotFound
		}
		return model.SystemSetting{}, fmt.Errorf("system_repository.LockForUpdate: %w", err)
	}
	return s, nil
}

// MarkInstalled sets installed=true (only if currently false) and writes related fields.
// Returns rowsAffected to allow caller to detect race conditions.
func (r *SystemRepository) MarkInstalled(ctx context.Context, tx *gorm.DB, siteName string, t time.Time) (int64, error) {
	if tx == nil {
		return 0, fmt.Errorf("system_repository.MarkInstalled: tx is required")
	}

	res := tx.WithContext(ctx).
		Model(&model.SystemSetting{}).
		Where("id = ? AND installed = ?", 1, false).
		Updates(map[string]any{
			"installed":    true,
			"site_name":    siteName,
			"installed_at": &t,
		})
	if res.Error != nil {
		return 0, fmt.Errorf("system_repository.MarkInstalled: %w", res.Error)
	}
	return res.RowsAffected, nil
}
