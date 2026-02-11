package dto

import (
	"KaldalisCMS/internal/core/entity"
	"time"
)

type MediaAssetResponse struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	OwnerUserID  uint      `json:"owner_user_id"`
	OriginalName string    `json:"original_name"`
	StoredName   string    `json:"stored_name"`
	Ext          string    `json:"ext"`
	MimeType     string    `json:"mime_type"`
	SizeBytes    int64     `json:"size_bytes"`
	Storage      string    `json:"storage"`
	ObjectKey    string    `json:"object_key"`
	Url          string    `json:"url"`
	Width        *int      `json:"width"`
	Height       *int      `json:"height"`
}

type MediaListResponse struct {
	Items    []MediaAssetResponse `json:"items"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
}

func ToMediaAssetResponse(a entity.MediaAsset) MediaAssetResponse {
	return MediaAssetResponse{
		ID:           a.ID,
		CreatedAt:    a.CreatedAt,
		OwnerUserID:  a.OwnerUserID,
		OriginalName: a.OriginalName,
		StoredName:   a.StoredName,
		Ext:          a.Ext,
		MimeType:     a.MimeType,
		SizeBytes:    a.SizeBytes,
		Storage:      a.Storage,
		ObjectKey:    a.ObjectKey,
		Url:          a.Url,
		Width:        a.Width,
		Height:       a.Height,
	}
}

func ToMediaAssetResponses(items []entity.MediaAsset) []MediaAssetResponse {
	out := make([]MediaAssetResponse, 0, len(items))
	for _, it := range items {
		out = append(out, ToMediaAssetResponse(it))
	}
	return out
}
