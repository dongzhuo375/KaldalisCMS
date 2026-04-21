package entity

import (
	"testing"
	"time"
)

func TestPost_CheckValidity(t *testing.T) {
	t.Run("empty title rejected", func(t *testing.T) {
		p := &Post{}
		if err := p.CheckValidity(); err == nil {
			t.Fatal("want error for empty title")
		}
	})
	t.Run("non-empty title ok", func(t *testing.T) {
		p := &Post{Title: "x"}
		if err := p.CheckValidity(); err != nil {
			t.Fatalf("unexpected: %v", err)
		}
	})
}

func TestPost_Publish(t *testing.T) {
	t.Run("from draft sets status and updated_at", func(t *testing.T) {
		p := &Post{Title: "x", Status: StatusDraft}
		before := time.Now().Add(-time.Second)
		if err := p.Publish(); err != nil {
			t.Fatal(err)
		}
		if p.Status != StatusPublished {
			t.Fatalf("status not updated: %d", p.Status)
		}
		if !p.UpdatedAt.After(before) {
			t.Fatalf("updated_at not bumped: %v", p.UpdatedAt)
		}
	})
	t.Run("already published rejected", func(t *testing.T) {
		p := &Post{Title: "x", Status: StatusPublished}
		if err := p.Publish(); err == nil {
			t.Fatal("want error for already published")
		}
	})
	t.Run("validity failure propagates", func(t *testing.T) {
		p := &Post{Status: StatusDraft}
		if err := p.Publish(); err == nil {
			t.Fatal("want error from empty title")
		}
	})
}

func TestPost_Draft(t *testing.T) {
	t.Run("already draft rejected", func(t *testing.T) {
		p := &Post{Status: StatusDraft}
		if err := p.Draft(); err == nil {
			t.Fatal("want error for already draft")
		}
	})
	t.Run("from published to draft", func(t *testing.T) {
		p := &Post{Status: StatusPublished}
		before := time.Now().Add(-time.Second)
		if err := p.Draft(); err != nil {
			t.Fatal(err)
		}
		if p.Status != StatusDraft {
			t.Fatalf("status: %d", p.Status)
		}
		if !p.UpdatedAt.After(before) {
			t.Fatalf("updated_at not bumped: %v", p.UpdatedAt)
		}
	})
}
