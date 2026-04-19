package dto

import (
	"testing"
	"time"

	"KaldalisCMS/internal/core/entity"
)

func TestCreatePostRequest_ToEntity(t *testing.T) {
	catID := uint(3)
	req := &CreatePostRequest{
		Title:      "hello",
		Content:    "body",
		Cover:      "/c.png",
		CategoryID: &catID,
		Tags:       []uint{1, 2},
	}
	got := req.ToEntity()
	if got.Title != "hello" || got.Content != "body" || got.Cover != "/c.png" {
		t.Fatalf("fields: %+v", got)
	}
	if got.Status != entity.StatusDraft {
		t.Fatalf("default status not draft: %d", got.Status)
	}
	if got.CategoryID == nil || *got.CategoryID != 3 {
		t.Fatalf("category lost: %+v", got.CategoryID)
	}
	if len(got.Tags) != 2 || got.Tags[0].ID != 1 || got.Tags[1].ID != 2 {
		t.Fatalf("tags: %+v", got.Tags)
	}
}

func TestCreatePostRequest_ToEntity_NilTags(t *testing.T) {
	req := &CreatePostRequest{Title: "x"}
	got := req.ToEntity()
	if got.Tags != nil {
		t.Fatalf("tags should stay nil when omitted, got %+v", got.Tags)
	}
}

func TestUpdatePostRequest_ToPatch_OnlySetFields(t *testing.T) {
	title := "new title"
	req := &UpdatePostRequest{Title: &title}
	patch := req.ToPatch()
	if patch.Title == nil || *patch.Title != "new title" {
		t.Fatalf("title: %+v", patch.Title)
	}
	if patch.Content != nil || patch.Cover != nil || patch.CategoryID != nil {
		t.Fatal("nil fields should remain nil in patch")
	}
	if patch.Tags != nil {
		t.Fatal("tags nil means 'do not touch'")
	}
}

func TestUpdatePostRequest_ToPatch_EmptyTagsMeansClear(t *testing.T) {
	// Per comment in PostPatch: nil = ignore, empty slice = replace with empty.
	req := &UpdatePostRequest{Tags: []uint{}}
	patch := req.ToPatch()
	if patch.Tags == nil {
		t.Fatal("empty tag slice must translate to non-nil empty (meaning 'clear')")
	}
	if len(patch.Tags) != 0 {
		t.Fatalf("len: %d", len(patch.Tags))
	}
}

func TestToPostResponse_Nil(t *testing.T) {
	if got := ToPostResponse(nil); got != nil {
		t.Fatalf("want nil, got %+v", got)
	}
}

func TestToPostResponse_FullMapping(t *testing.T) {
	catID := uint(9)
	p := &entity.Post{
		ID:        1,
		Title:     "t",
		Slug:      "t",
		Content:   "c",
		Cover:     "/c.png",
		Status:    entity.StatusPublished,
		CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
		Author:    entity.User{ID: 5, Username: "alice"},
		CategoryID: &catID,
		Category:  entity.Category{ID: 9, Name: "tech"},
		Tags:      []entity.Tag{{ID: 10, Name: "go"}},
	}
	got := ToPostResponse(p)
	if got.ID != 1 || got.Title != "t" || got.Status != entity.StatusPublished {
		t.Fatalf("scalar: %+v", got)
	}
	if got.Author.Username != "alice" {
		t.Fatalf("author: %+v", got.Author)
	}
	if got.Category == nil || got.Category.Name != "tech" {
		t.Fatalf("category: %+v", got.Category)
	}
	if len(got.Tags) != 1 || got.Tags[0].Name != "go" {
		t.Fatalf("tags: %+v", got.Tags)
	}
	if got.CreatedAt != "2026-01-01T00:00:00Z" {
		t.Fatalf("time format: %q", got.CreatedAt)
	}
}

func TestToPostResponse_NoCategoryNoTags(t *testing.T) {
	p := &entity.Post{ID: 1, Title: "t"}
	got := ToPostResponse(p)
	if got.Category != nil {
		t.Fatalf("category should be nil when CategoryID nil: %+v", got.Category)
	}
	if got.Tags != nil {
		t.Fatalf("tags should be nil when empty: %+v", got.Tags)
	}
}

func TestToPostListResponse_Empty(t *testing.T) {
	got := ToPostListResponse(nil)
	if got == nil {
		t.Fatal("want empty slice, got nil (breaks JSON contract)")
	}
	if len(got) != 0 {
		t.Fatalf("len: %d", len(got))
	}
}

func TestToPostListResponse_Many(t *testing.T) {
	got := ToPostListResponse([]entity.Post{{ID: 1, Title: "a"}, {ID: 2, Title: "b"}})
	if len(got) != 2 || got[0].ID != 1 || got[1].ID != 2 {
		t.Fatalf("unexpected: %+v", got)
	}
}
