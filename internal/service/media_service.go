package service

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type MediaConfig struct {
	UploadDir        string
	MaxUploadSizeMB  int64
	PublicBaseURL    string
	MaxFilenameBytes int
}

type MediaService struct {
	repo core.MediaRepository
	cfg  MediaConfig
}

func NewMediaService(repo core.MediaRepository, cfg MediaConfig) *MediaService {
	if cfg.MaxUploadSizeMB <= 0 {
		cfg.MaxUploadSizeMB = 50
	}
	if cfg.MaxFilenameBytes <= 0 {
		cfg.MaxFilenameBytes = 180
	}
	if cfg.UploadDir == "" {
		cfg.UploadDir = filepath.FromSlash("./data/uploads")
	}
	return &MediaService{repo: repo, cfg: cfg}
}

var (
	ErrUploadTooLarge   = errors.New("upload too large")
	ErrUnsupportedType  = errors.New("unsupported file type")
	ErrAssetReferenced  = errors.New("asset is referenced by posts")
	ErrInvalidAssetName = errors.New("invalid asset name")
)

// CreateAssetFromUpload persists metadata and stores file under:
// {upload_dir}/a/{assetID}/{stored_name}
// Public URL:
// {public_base_url}/media/a/{assetID}/{stored_name}  (public_base_url may be empty)
func (s *MediaService) CreateAssetFromUpload(ctx context.Context, ownerUserID uint, fileHeader *multipart.FileHeader) (entity.MediaAsset, error) {
	if fileHeader == nil {
		return entity.MediaAsset{}, fmt.Errorf("media_service.CreateAssetFromUpload: file is nil")
	}

	maxBytes := s.cfg.MaxUploadSizeMB * 1024 * 1024
	if maxBytes > 0 && fileHeader.Size > maxBytes {
		return entity.MediaAsset{}, ErrUploadTooLarge
	}

	origName := fileHeader.Filename
	storedName, ext, err := sanitizeFilename(origName, s.cfg.MaxFilenameBytes)
	if err != nil {
		return entity.MediaAsset{}, err
	}

	f, err := fileHeader.Open()
	if err != nil {
		return entity.MediaAsset{}, fmt.Errorf("media_service.Open: %w", err)
	}
	defer f.Close()

	// Detect mime type from first 512 bytes.
	sniff := make([]byte, 512)
	n, _ := io.ReadFull(f, sniff)
	sniff = sniff[:n]
	mimeType := http.DetectContentType(sniff)

	if !isAllowedMime(mimeType) {
		return entity.MediaAsset{}, ErrUnsupportedType
	}

	// Reset reader: reopen (multipart.File does not necessarily support Seek)
	_ = f.Close()
	f, err = fileHeader.Open()
	if err != nil {
		return entity.MediaAsset{}, fmt.Errorf("media_service.reopen: %w", err)
	}
	defer f.Close()

	asset := entity.MediaAsset{
		OwnerUserID:  ownerUserID,
		OriginalName: origName,
		StoredName:   storedName,
		Ext:          ext,
		MimeType:     mimeType,
		SizeBytes:    fileHeader.Size,
		Storage:      "local",
	}

	// Create row first to get asset.ID.
	if err := s.repo.Create(ctx, &asset); err != nil {
		return entity.MediaAsset{}, err
	}

	objectKey := filepath.ToSlash(filepath.Join("a", fmt.Sprintf("%d", asset.ID), storedName))
	asset.ObjectKey = objectKey
	asset.Url = joinPublicURL(s.cfg.PublicBaseURL, "/media/"+objectKey)

	// Best-effort image config (only for images)
	if strings.HasPrefix(strings.ToLower(mimeType), "image/") {
		if w, h := tryReadImageSize(fileHeader); w != nil && h != nil {
			asset.Width = w
			asset.Height = h
		}
	}

	// Persist ObjectKey/Url/Width/Height
	if err := s.repo.UpdateAssetFields(ctx, asset.ID, map[string]any{"object_key": asset.ObjectKey, "url": asset.Url, "width": asset.Width, "height": asset.Height}); err != nil {
		return entity.MediaAsset{}, fmt.Errorf("media_service.update_metadata: %w", err)
	}

	// Ensure directory and write file.
	absPath := filepath.Join(s.cfg.UploadDir, filepath.FromSlash(objectKey))
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return entity.MediaAsset{}, fmt.Errorf("media_service.MkdirAll: %w", err)
	}

	out, err := os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return entity.MediaAsset{}, fmt.Errorf("media_service.OpenFile: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, f); err != nil {
		return entity.MediaAsset{}, fmt.Errorf("media_service.Copy: %w", err)
	}

	return asset, nil
}

func (s *MediaService) List(ctx context.Context, requesterRole string, requesterUserID uint, page, pageSize int, q string) ([]entity.MediaAsset, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var owner *uint
	if requesterRole != "admin" {
		owner = &requesterUserID
	}
	return s.repo.List(ctx, owner, offset, pageSize, q)
}

// DeleteAs deletes an asset with permission checks:
// - admin can delete any asset (still subject to reference hard restriction)
// - non-admin can only delete assets they own
func (s *MediaService) DeleteAs(ctx context.Context, requesterRole string, requesterUserID uint, assetID uint) error {
	asset, err := s.repo.GetByID(ctx, assetID)
	if err != nil {
		return err
	}
	if requesterRole != "admin" && asset.OwnerUserID != requesterUserID {
		return core.ErrPermission
	}
	return s.deleteByAsset(ctx, asset)
}

// Delete keeps legacy behavior (no owner check) and is treated as a privileged internal operation.
func (s *MediaService) Delete(ctx context.Context, assetID uint) error {
	asset, err := s.repo.GetByID(ctx, assetID)
	if err != nil {
		return err
	}
	return s.deleteByAsset(ctx, asset)
}

func (s *MediaService) deleteByAsset(ctx context.Context, asset entity.MediaAsset) error {
	cnt, err := s.repo.CountReferences(ctx, asset.ID)
	if err != nil {
		return err
	}
	if cnt > 0 {
		return ErrAssetReferenced
	}

	if err := s.repo.Delete(ctx, asset.ID); err != nil {
		return err
	}

	// Best-effort remove physical file + dir.
	absPath := filepath.Join(s.cfg.UploadDir, filepath.FromSlash(asset.ObjectKey))
	_ = os.Remove(absPath)
	_ = os.Remove(filepath.Dir(absPath))
	return nil
}

func (s *MediaService) ListPostMedia(ctx context.Context, postID uint, purpose *string) ([]entity.MediaAsset, error) {
	return s.repo.ListPostMedia(ctx, postID, purpose)
}

// SyncPostReferences parses markdown content and cover URL to update post_assets mappings.
func (s *MediaService) SyncPostReferences(ctx context.Context, postID uint, content string, cover string) error {
	contentIDs := extractAssetIDsFromMarkdown(content)
	if err := s.repo.UpsertPostReferences(ctx, postID, "content", contentIDs); err != nil {
		return err
	}

	coverID := extractAssetIDFromMediaURL(cover)
	coverIDs := []uint{}
	if coverID != 0 {
		coverIDs = []uint{coverID}
	}
	if err := s.repo.UpsertPostReferences(ctx, postID, "cover", coverIDs); err != nil {
		return err
	}
	return nil
}

// --- helpers ---

func joinPublicURL(base, path string) string {
	if base == "" {
		return path
	}
	return strings.TrimRight(base, "/") + path
}

var invalidNameChars = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func sanitizeFilename(name string, maxBytes int) (stored string, ext string, err error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", "", ErrInvalidAssetName
	}
	base := filepath.Base(name)
	ext = strings.ToLower(filepath.Ext(base))
	stem := strings.TrimSuffix(base, ext)
	stem = invalidNameChars.ReplaceAllString(stem, "-")
	stem = strings.Trim(stem, "-._")
	if stem == "" {
		stem = uuid.New().String()
	}
	if ext == "" {
		ext = ".bin"
	}
	stored = stem + ext
	if maxBytes > 0 && len([]byte(stored)) > maxBytes {
		maxStem := maxBytes - len([]byte(ext))
		if maxStem <= 8 {
			stored = uuid.New().String() + ext
		} else {
			b := []byte(stem)
			if maxStem > len(b) {
				maxStem = len(b)
			}
			stored = string(b[:maxStem]) + ext
		}
	}
	return stored, ext, nil
}

func isAllowedMime(mimeType string) bool {
	m := strings.ToLower(mimeType)
	if strings.HasPrefix(m, "image/") || strings.HasPrefix(m, "video/") || strings.HasPrefix(m, "audio/") {
		return true
	}
	switch m {
	case "application/pdf":
		return true
	case "application/zip", "application/x-zip-compressed":
		return true
	case "application/x-rar-compressed", "application/vnd.rar":
		return true
	case "application/x-7z-compressed":
		return true
	default:
		return false
	}
}

var reAssetURL = regexp.MustCompile(`/media/a/(\d+)/[^)\s]+`)

func extractAssetIDsFromMarkdown(md string) []uint {
	matches := reAssetURL.FindAllStringSubmatch(md, -1)
	if len(matches) == 0 {
		return nil
	}
	seen := make(map[uint]struct{}, len(matches))
	out := make([]uint, 0, len(matches))
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		id := parseUint(m[1])
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func extractAssetIDFromMediaURL(u string) uint {
	if u == "" {
		return 0
	}
	m := reAssetURL.FindStringSubmatch(u)
	if len(m) < 2 {
		return 0
	}
	return parseUint(m[1])
}

func parseUint(s string) uint {
	var v uint64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0
		}
		v = v*10 + uint64(c-'0')
		if v > uint64(^uint(0)) {
			return 0
		}
	}
	return uint(v)
}

func tryReadImageSize(fileHeader *multipart.FileHeader) (*int, *int) {
	f, err := fileHeader.Open()
	if err != nil {
		return nil, nil
	}
	defer f.Close()

	cfg, _, err := decodeImageConfig(f)
	if err != nil {
		return nil, nil
	}
	w, h := cfg.Width, cfg.Height
	return &w, &h
}
