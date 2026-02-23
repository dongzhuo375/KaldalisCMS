package entity

import "time"

// MediaStatus represents the state of a media asset upload.
type MediaStatus int

const (
	MediaStatusPending  MediaStatus = 0 // Initial state, record created but file not yet confirmed on disk.
	MediaStatusUploaded MediaStatus = 1 // File successfully written and verified.
	MediaStatusFailed   MediaStatus = 2 // Upload failed or file write error.
)

// MediaAsset is the domain entity for an uploaded media resource.
// Files are served publicly via: /media/a/{id}/{stored_name}
type MediaAsset struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	OwnerUserID uint

	OriginalName string
	StoredName   string
	Ext          string
	MimeType     string
	SizeBytes    int64

	SHA256    string
	Storage   string
	ObjectKey string
	Url       string

	Width  *int
	Height *int

	// Status tracks the lifecycle of the asset (PENDING -> UPLOADED / FAILED)
	Status MediaStatus
}
