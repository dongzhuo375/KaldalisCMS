package service

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
)

// --- pure helpers ---

func TestSanitizeFilename(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		maxBytes int
		wantExt string
		wantErr bool
		check   func(t *testing.T, stored string)
	}{
		{
			// Ext is lowercased but stem keeps original case.
			// TrimSuffix is case-sensitive, so "photo.JPG" - ".jpg" is a no-op and the original case survives in the stem.
			name: "lowercase extension",
			input: "photo.jpg", maxBytes: 100, wantExt: ".jpg",
			check: func(t *testing.T, s string) {
				if s != "photo.jpg" {
					t.Fatalf("stored: %q", s)
				}
			},
		},
		{
			name: "spaces and special chars replaced",
			input: "my file (1).png", maxBytes: 100, wantExt: ".png",
			check: func(t *testing.T, s string) {
				if strings.ContainsAny(s, " ()") {
					t.Fatalf("special chars not replaced: %q", s)
				}
			},
		},
		{
			name: "no extension defaults to .bin",
			input: "README", maxBytes: 100, wantExt: ".bin",
			check: func(t *testing.T, s string) {
				if !strings.HasSuffix(s, ".bin") {
					t.Fatalf("ext fallback: %q", s)
				}
			},
		},
		{
			name:    "empty rejected",
			input:   "   ",
			maxBytes: 100,
			wantErr: true,
		},
		{
			name: "overlong truncated",
			input: strings.Repeat("a", 500) + ".png", maxBytes: 30, wantExt: ".png",
			check: func(t *testing.T, s string) {
				if len(s) > 30 {
					t.Fatalf("not truncated: len=%d", len(s))
				}
				if !strings.HasSuffix(s, ".png") {
					t.Fatalf("ext lost: %q", s)
				}
			},
		},
		{
			name: "all invalid chars → uuid stem",
			input: "@@@.png", maxBytes: 100, wantExt: ".png",
			check: func(t *testing.T, s string) {
				if !strings.HasSuffix(s, ".png") {
					t.Fatalf("ext: %q", s)
				}
				if s == ".png" {
					t.Fatalf("empty stem not replaced: %q", s)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stored, ext, err := sanitizeFilename(tc.input, tc.maxBytes)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("want error, got stored=%q", stored)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if ext != tc.wantExt {
				t.Fatalf("ext: got %q want %q", ext, tc.wantExt)
			}
			if tc.check != nil {
				tc.check(t, stored)
			}
		})
	}
}

func TestIsAllowedMime(t *testing.T) {
	allowed := []string{
		"image/png", "image/jpeg", "video/mp4", "audio/mpeg",
		"application/pdf", "application/zip",
		"application/x-rar-compressed", "application/vnd.rar",
		"application/x-7z-compressed",
	}
	for _, m := range allowed {
		if !isAllowedMime(m) {
			t.Errorf("should allow %q", m)
		}
	}

	denied := []string{
		"application/octet-stream", "text/html", "application/javascript",
		"", "unknown/type",
	}
	for _, m := range denied {
		if isAllowedMime(m) {
			t.Errorf("should deny %q", m)
		}
	}
}

func TestExtractAssetIDsFromMarkdown(t *testing.T) {
	md := `
		![alt](/media/a/1/a.png)
		some text
		![dup](/media/a/1/b.jpg)
		[link](/media/a/42/c.pdf)
		not a media: /other/a/99/x.png
	`
	got := extractAssetIDsFromMarkdown(md)
	want := []uint{1, 42}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}

	if got := extractAssetIDsFromMarkdown(""); got != nil {
		t.Fatalf("empty markdown should return nil, got %v", got)
	}
}

func TestExtractAssetIDFromMediaURL(t *testing.T) {
	if id := extractAssetIDFromMediaURL(""); id != 0 {
		t.Fatalf("empty: %d", id)
	}
	if id := extractAssetIDFromMediaURL("/media/a/7/cover.png"); id != 7 {
		t.Fatalf("valid: %d", id)
	}
	if id := extractAssetIDFromMediaURL("/not-media/a/7/x"); id != 0 {
		t.Fatalf("non-media path should be 0, got %d", id)
	}
}

func TestParseUint(t *testing.T) {
	if parseUint("") != 0 {
		t.Fatal("empty → 0")
	}
	if parseUint("123") != 123 {
		t.Fatal("decimal")
	}
	if parseUint("12a") != 0 {
		t.Fatal("non-digit → 0")
	}
}

func TestJoinPublicURL(t *testing.T) {
	if got := joinPublicURL("", "/media/a/1/x.png"); got != "/media/a/1/x.png" {
		t.Fatalf("empty base: %q", got)
	}
	if got := joinPublicURL("https://cdn.example.com", "/media/a/1/x.png"); got != "https://cdn.example.com/media/a/1/x.png" {
		t.Fatalf("with base: %q", got)
	}
	if got := joinPublicURL("https://cdn.example.com/", "/x"); got != "https://cdn.example.com/x" {
		t.Fatalf("trailing slash not trimmed: %q", got)
	}
}

// --- NewMediaService defaults ---

func TestNewMediaService_Defaults(t *testing.T) {
	s := NewMediaService(nil, MediaConfig{})
	if s.cfg.MaxUploadSizeMB != 50 {
		t.Fatalf("max size default: %d", s.cfg.MaxUploadSizeMB)
	}
	if s.cfg.MaxFilenameBytes != 180 {
		t.Fatalf("filename default: %d", s.cfg.MaxFilenameBytes)
	}
	if s.cfg.UploadDir == "" {
		t.Fatal("upload dir default missing")
	}
}

// --- List (repo-wrapping with pagination & owner scoping) ---

type fakeMediaRepoForList struct {
	fakeMediaRepoNoOp
	listFn func(ctx context.Context, owner *uint, offset, limit int, q string) ([]entity.MediaAsset, int64, error)
}

func (f *fakeMediaRepoForList) List(ctx context.Context, owner *uint, offset, limit int, q string) ([]entity.MediaAsset, int64, error) {
	return f.listFn(ctx, owner, offset, limit, q)
}

func TestMediaService_List_AdminSeesAll(t *testing.T) {
	repo := &fakeMediaRepoForList{
		listFn: func(ctx context.Context, owner *uint, offset, limit int, q string) ([]entity.MediaAsset, int64, error) {
			if owner != nil {
				t.Fatalf("admin should not be scoped by owner, got %d", *owner)
			}
			if limit != 20 || offset != 0 {
				t.Fatalf("default pagination: offset=%d limit=%d", offset, limit)
			}
			return []entity.MediaAsset{{ID: 1}}, 1, nil
		},
	}
	svc := NewMediaService(repo, MediaConfig{})
	assets, total, err := svc.List(context.Background(), "admin", 9, 0, 0, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(assets) != 1 || total != 1 {
		t.Fatalf("got %d/%d", len(assets), total)
	}
}

func TestMediaService_List_UserScopedToOwn(t *testing.T) {
	repo := &fakeMediaRepoForList{
		listFn: func(ctx context.Context, owner *uint, offset, limit int, q string) ([]entity.MediaAsset, int64, error) {
			if owner == nil || *owner != 5 {
				t.Fatalf("owner scope: %+v", owner)
			}
			return nil, 0, nil
		},
	}
	_, _, err := NewMediaService(repo, MediaConfig{}).List(context.Background(), "editor", 5, 1, 10, "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMediaService_List_PageSizeClamped(t *testing.T) {
	repo := &fakeMediaRepoForList{
		listFn: func(ctx context.Context, owner *uint, offset, limit int, q string) ([]entity.MediaAsset, int64, error) {
			if limit != 100 {
				t.Fatalf("limit should be clamped to 100, got %d", limit)
			}
			return nil, 0, nil
		},
	}
	_, _, err := NewMediaService(repo, MediaConfig{}).List(context.Background(), "admin", 9, 1, 9999, "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMediaService_List_RepoErrorNormalized(t *testing.T) {
	repo := &fakeMediaRepoForList{
		listFn: func(ctx context.Context, owner *uint, offset, limit int, q string) ([]entity.MediaAsset, int64, error) {
			return nil, 0, errors.New("db down")
		},
	}
	_, _, err := NewMediaService(repo, MediaConfig{}).List(context.Background(), "admin", 9, 1, 10, "")
	if !errors.Is(err, core.ErrInternalError) {
		t.Fatalf("want ErrInternalError, got %v", err)
	}
}

// --- fakeMediaRepoNoOp: satisfies core.MediaRepository; all methods panic unless overridden.
// Use embedded composition to override only the method(s) a test cares about.

type fakeMediaRepoNoOp struct{}

func (fakeMediaRepoNoOp) Create(ctx context.Context, asset *entity.MediaAsset) error { panic("not impl") }
func (fakeMediaRepoNoOp) GetByID(ctx context.Context, id uint) (entity.MediaAsset, error) {
	panic("not impl")
}
func (fakeMediaRepoNoOp) List(ctx context.Context, owner *uint, offset, limit int, q string) ([]entity.MediaAsset, int64, error) {
	panic("not impl")
}
func (fakeMediaRepoNoOp) Delete(ctx context.Context, id uint) error { panic("not impl") }
func (fakeMediaRepoNoOp) CountReferences(ctx context.Context, assetID uint) (int64, error) {
	panic("not impl")
}
func (fakeMediaRepoNoOp) UpsertPostReferences(ctx context.Context, postID uint, purpose string, assetIDs []uint) error {
	panic("not impl")
}
func (fakeMediaRepoNoOp) ListPostMedia(ctx context.Context, postID uint, purpose *string) ([]entity.MediaAsset, error) {
	panic("not impl")
}
func (fakeMediaRepoNoOp) UpdateAssetFields(ctx context.Context, assetID uint, fields map[string]any) error {
	panic("not impl")
}
func (fakeMediaRepoNoOp) UpdateStatus(ctx context.Context, id uint, status entity.MediaStatus) error {
	panic("not impl")
}
func (fakeMediaRepoNoOp) ListPendingOlderThan(ctx context.Context, cutoff time.Time, limit int) ([]entity.MediaAsset, error) {
	panic("not impl")
}
func (fakeMediaRepoNoOp) ListSoftDeletedOlderThan(ctx context.Context, cutoff time.Time, limit int) ([]entity.MediaAsset, error) {
	panic("not impl")
}
func (fakeMediaRepoNoOp) DeletePhysical(ctx context.Context, id uint) error { panic("not impl") }
