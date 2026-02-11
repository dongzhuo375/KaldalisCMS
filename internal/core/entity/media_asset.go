package entity

import "time"

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
}
