package model

import (
	"time"

	"gorm.io/gorm"
)

// PostAsset maps which media assets are referenced by a post.
// Purpose:
//   - content: referenced from post.content markdown
//   - cover: referenced from post.cover
//
// Unique constraint prevents duplicates.
type PostAsset struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	PostID  uint   `gorm:"not null;index;uniqueIndex:idx_post_asset" json:"post_id"`
	AssetID uint   `gorm:"not null;index;uniqueIndex:idx_post_asset" json:"asset_id"`
	Purpose string `gorm:"not null;default:'content';uniqueIndex:idx_post_asset" json:"purpose"`
}
