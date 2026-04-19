package dto

import (
	"testing"

	"KaldalisCMS/internal/core/entity"
)

func TestCreateTagRequest_ToEntity(t *testing.T) {
	r := &CreateTagRequest{Name: "go"}
	got := r.ToEntity()
	if got.Name != "go" {
		t.Fatalf("name: %q", got.Name)
	}
}

func TestUpdateTagRequest_ToEntity(t *testing.T) {
	name := "new"
	r := &UpdateTagRequest{Name: &name}
	got := r.ToEntity()
	if got.Name != "new" {
		t.Fatalf("name: %q", got.Name)
	}

	r2 := &UpdateTagRequest{}
	got2 := r2.ToEntity()
	if got2.Name != "" {
		t.Fatalf("nil name should leave empty: %q", got2.Name)
	}
}

func TestToTagResponse(t *testing.T) {
	got := ToTagResponse(entity.Tag{ID: 1, Name: "go"})
	if got.ID != 1 || got.Name != "go" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestToTagListResponse_Empty(t *testing.T) {
	got := ToTagListResponse(nil)
	if got == nil {
		t.Fatal("want empty slice not nil (JSON contract)")
	}
	if len(got) != 0 {
		t.Fatalf("len: %d", len(got))
	}
}

func TestToTagListResponse_Many(t *testing.T) {
	got := ToTagListResponse([]entity.Tag{{ID: 1}, {ID: 2}})
	if len(got) != 2 || got[0].ID != 1 || got[1].ID != 2 {
		t.Fatalf("unexpected: %+v", got)
	}
}
