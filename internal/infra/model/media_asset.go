package model

import (
	"time"

	"gorm.io/gorm"
)

// MediaAsset stores metadata for an uploaded file.
// Files are served publicly via: /media/a/{id}/{stored_name}
// and stored on disk under: {upload_dir}/a/{id}/{stored_name}
//
// NOTE: Referencing is maintained via PostAsset.
type MediaAsset struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OwnerUserID uint `gorm:"not null;index" json:"owner_user_id"`

	OriginalName string `gorm:"not null" json:"original_name"`
	StoredName   string `gorm:"not null" json:"stored_name"`
	Ext          string `gorm:"not null" json:"ext"`
	MimeType     string `gorm:"not null" json:"mime_type"`
	SizeBytes    int64  `gorm:"not null" json:"size_bytes"`

	// SHA256 is optional today and can be used later for dedupe/verification.
	SHA256 string `gorm:"" json:"sha256"`

	// Storage is reserved for future backends (s3/minio). For now: "local".
	Storage string `gorm:"not null;default:'local'" json:"storage"`

	// ObjectKey is the relative key inside UploadDir. For local storage it matches disk path.
	// Example: a/123/photo.png
	ObjectKey string `gorm:"not null;uniqueIndex" json:"object_key"`

	// Url is the public URL used by Markdown content.
	// Example: /media/a/123/photo.png or https://cdn.example.com/media/a/123/photo.png
	Url string `gorm:"not null" json:"url"`

	// Optional image metadata.
	Width  *int `json:"width"`
	Height *int `json:"height"`
}
