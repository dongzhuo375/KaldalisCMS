package service

import (
	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
	repository "KaldalisCMS/internal/infra/repository/postgres"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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
	ErrUploadTooLarge   = fmt.Errorf("%w: upload too large", core.ErrInvalidInput)
	ErrUnsupportedType  = fmt.Errorf("%w: unsupported file type", core.ErrInvalidInput)
	ErrAssetReferenced  = fmt.Errorf("%w: asset is referenced by posts", core.ErrConflict)
	ErrInvalidAssetName = fmt.Errorf("%w: invalid asset name", core.ErrInvalidInput)
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

	// --- State Machine Step 1: PENDING ---
	// Insert DB record first to get ID. Status defaults to PENDING (0).
	asset := entity.MediaAsset{
		OwnerUserID:  ownerUserID,
		OriginalName: origName,
		StoredName:   storedName,
		Ext:          ext,
		MimeType:     mimeType,
		SizeBytes:    fileHeader.Size,
		Storage:      "local",
		Status:       entity.MediaStatusPending,
	}

	if err := s.repo.Create(ctx, &asset); err != nil {
		return entity.MediaAsset{}, err
	}

	// Calculate paths
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

	// --- State Machine Step 2: WRITE FILE ---
	absPath := filepath.Join(s.cfg.UploadDir, filepath.FromSlash(objectKey))
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		// Write failed -> Mark as FAILED
		_ = s.repo.UpdateStatus(ctx, asset.ID, entity.MediaStatusFailed)
		return entity.MediaAsset{}, fmt.Errorf("media_service.MkdirAll: %w", err)
	}

	out, err := os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		// Write failed -> Mark as FAILED
		_ = s.repo.UpdateStatus(ctx, asset.ID, entity.MediaStatusFailed)
		return entity.MediaAsset{}, fmt.Errorf("media_service.OpenFile: %w", err)
	}
	// We must close explicitely to ensure flush before DB update
	// so we don't rely on defer alone for the success path
	copyErr := func() error {
		defer out.Close()
		if _, err := io.Copy(out, f); err != nil {
			return err
		}
		return nil
	}()

	if copyErr != nil {
		// Copy failed -> Mark as FAILED
		// Try to clean up partial file
		_ = os.Remove(absPath)
		_ = s.repo.UpdateStatus(ctx, asset.ID, entity.MediaStatusFailed)
		return entity.MediaAsset{}, fmt.Errorf("media_service.Copy: %w", copyErr)
	}

	// --- State Machine Step 3: UPLOADED ---
	// File is safe on disk. Update metadata and flip status to UPLOADED.
	updates := map[string]any{
		"object_key": asset.ObjectKey,
		"url":        asset.Url,
		"width":      asset.Width,
		"height":     asset.Height,
		"status":     int(entity.MediaStatusUploaded),
	}

	if err := s.repo.UpdateAssetFields(ctx, asset.ID, updates); err != nil {
		// DB update failed. This is the "Inconsistent" state (File ok, DB pending).
		// We leave it as PENDING. The background cleanup job will see it's old and delete the file + record.
		// Alternatively, we could try to delete the file here, but let's rely on the cleanup job for robustness.
		return entity.MediaAsset{}, fmt.Errorf("media_service.update_metadata: %w", err)
	}

	asset.Status = entity.MediaStatusUploaded
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
		if errors.Is(err, repository.ErrMediaNotFound) {
			return core.ErrNotFound
		}
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
		if errors.Is(err, repository.ErrMediaNotFound) {
			return core.ErrNotFound
		}
		return err
	}
	return s.deleteByAsset(ctx, asset)
}

// CleanupStaleMedia scans for stale PENDING assets AND soft-deleted assets and removes them.
// Providing a robust "Eventual Consistency" guarantee.
func (s *MediaService) CleanupStaleMedia(ctx context.Context) error {
	// 1. Clean up Pending assets (older than 1 hour)
	cutoffPending := time.Now().Add(-1 * time.Hour)
	limit := 50 // process in batches

	staleAssets, err := s.repo.ListPendingOlderThan(ctx, cutoffPending, limit)
	if err != nil {
		fmt.Printf("[MediaCleanup] Failed to list pending assets: %v\n", err)
	} else if len(staleAssets) > 0 {
		fmt.Printf("[MediaCleanup] Found %d stale pending assets. Cleaning up...\n", len(staleAssets))
		for _, asset := range staleAssets {
			s.physicalDelete(ctx, asset)
		}
	}

	// 2. Clean up Soft-Deleted assets (older than 24 hours to allow accidental recovery if needed, or immediate if preferred)
	// For this requirement (Scheme A), we can clean them up relatively quickly, e.g., 1 hour or even 5 minutes.
	// Let's use 1 hour to be safe and consistent with Pending cleanup.
	cutoffDeleted := time.Now().Add(-1 * time.Hour)
	deletedAssets, err := s.repo.ListSoftDeletedOlderThan(ctx, cutoffDeleted, limit)
	if err != nil {
		return fmt.Errorf("failed to list soft-deleted assets: %w", err)
	}

	if len(deletedAssets) > 0 {
		fmt.Printf("[MediaCleanup] Found %d soft-deleted assets. Finalizing cleanup...\n", len(deletedAssets))
		for _, asset := range deletedAssets {
			s.physicalDelete(ctx, asset)
		}
	}

	return nil
}

func (s *MediaService) physicalDelete(ctx context.Context, asset entity.MediaAsset) {
	// 1. Delete physical file (if exists)
	objectKey := filepath.ToSlash(filepath.Join("a", fmt.Sprintf("%d", asset.ID), asset.StoredName))
	absPath := filepath.Join(s.cfg.UploadDir, filepath.FromSlash(objectKey))

	// Best effort remove file and dir
	// We check if file exists to give better logs, but Remove is idempotent-ish
	if _, err := os.Stat(absPath); err == nil {
		if err := os.Remove(absPath); err != nil {
			fmt.Printf("[MediaCleanup] Failed to remove file %s: %v\n", absPath, err)
			return // If file deletion fails (e.g. locked), retry later
		}
	}
	_ = os.Remove(filepath.Dir(absPath)) // try remove dir

	// 2. Delete DB record HARD
	if err := s.repo.DeletePhysical(ctx, asset.ID); err != nil {
		fmt.Printf("[MediaCleanup] Failed to hard delete asset ID %d: %v\n", asset.ID, err)
	} else {
		fmt.Printf("[MediaCleanup] Hard deleted asset ID %d\n", asset.ID)
	}
}

func (s *MediaService) deleteByAsset(ctx context.Context, asset entity.MediaAsset) error {
	cnt, err := s.repo.CountReferences(ctx, asset.ID)
	if err != nil {
		return err
	}
	if cnt > 0 {
		return ErrAssetReferenced
	}

	// Scheme A: Soft Delete Only
	// We rely on background job to cleanup the file and hard delete the record.
	if err := s.repo.Delete(ctx, asset.ID); err != nil {
		return err
	}

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

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return nil, nil
	}
	w, h := cfg.Width, cfg.Height
	return &w, &h
}
