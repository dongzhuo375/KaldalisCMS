package dto

import (
	"testing"
	"time"

	"KaldalisCMS/internal/core/entity"
)

func TestToMediaAssetResponse(t *testing.T) {
	w := 800
	h := 600
	created := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	a := entity.MediaAsset{
		ID:           5,
		CreatedAt:    created,
		OwnerUserID:  7,
		OriginalName: "a.png",
		StoredName:   "abc.png",
		Ext:          "png",
		MimeType:     "image/png",
		SizeBytes:    1024,
		Storage:      "local",
		ObjectKey:    "a/5/abc.png",
		Url:          "/media/a/5/abc.png",
		Width:        &w,
		Height:       &h,
		Status:       entity.MediaStatusUploaded,
	}
	got := ToMediaAssetResponse(a)
	if got.ID != 5 || got.Url != "/media/a/5/abc.png" {
		t.Fatalf("scalar: %+v", got)
	}
	if got.Width == nil || *got.Width != 800 {
		t.Fatalf("width: %+v", got.Width)
	}
	if got.Status != int(entity.MediaStatusUploaded) {
		t.Fatalf("status: %d", got.Status)
	}
}

func TestToMediaAssetResponses_Empty(t *testing.T) {
	got := ToMediaAssetResponses(nil)
	if got == nil {
		t.Fatal("want empty slice, got nil")
	}
	if len(got) != 0 {
		t.Fatalf("len: %d", len(got))
	}
}

func TestToMediaAssetResponses_Many(t *testing.T) {
	got := ToMediaAssetResponses([]entity.MediaAsset{{ID: 1}, {ID: 2}})
	if len(got) != 2 || got[0].ID != 1 || got[1].ID != 2 {
		t.Fatalf("unexpected: %+v", got)
	}
}
