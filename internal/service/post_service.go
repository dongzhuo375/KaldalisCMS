package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gosimple/slug"

	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
)

// PostService orchestrates the article publishing workflow.
// It deliberately exposes separate public and management entry points so callers do not
// have to understand which statuses are externally visible.
type PostService struct {
	repo       core.PostRepository
	authorizer core.PostAuthorizer
	// media is optional; when nil, reference sync is skipped.
	media *MediaService
}

func NewPostService(repo core.PostRepository, authorizer core.PostAuthorizer) *PostService {
	return &PostService{repo: repo, authorizer: authorizer}
}

// NewPostServiceWithMedia wires optional collaborators used by the post workflow.
// Authorization is delegated to the injected authorizer so service code never hard-codes role names.
func NewPostServiceWithMedia(repo core.PostRepository, media *MediaService, authorizer core.PostAuthorizer) *PostService {
	return &PostService{repo: repo, media: media, authorizer: authorizer}
}

// ListPublicPosts returns only published posts.
func (s *PostService) ListPublicPosts(ctx context.Context) ([]entity.Post, error) {
	posts, err := s.repo.GetPublished(ctx)
	if err != nil {
		return nil, normalizeServiceErrorWithOpMsg("post.list_public", "list published posts failed", err)
	}
	return posts, nil
}

// GetPublicPostByID returns a single published post for anonymous/public readers.
func (s *PostService) GetPublicPostByID(ctx context.Context, id uint) (entity.Post, error) {
	post, err := s.repo.GetPublishedByID(ctx, id)
	if err != nil {
		return entity.Post{}, normalizeServiceErrorWithOpMsg("post.get_public_by_id", "get published post by id failed", err)
	}
	return post, nil
}

// ListAdminPosts returns the management view of posts for the acting user.
func (s *PostService) ListAdminPosts(ctx context.Context, actorUserID uint, actorRole string) ([]entity.Post, error) {
	canListAny, err := s.hasPostPermission(ctx, actorRole, core.PostPermissionListAnyPost)
	if err != nil {
		return nil, err
	}
	if canListAny {
		posts, err := s.repo.GetAll(ctx)
		if err != nil {
			return nil, normalizeServiceErrorWithOpMsg("post.list_admin_all", "list all posts in admin scope failed", err)
		}
		return posts, nil
	}

	if actorUserID == 0 {
		return nil, core.ErrPermission
	}
	if err := s.authorizePostAction(ctx, actorRole, core.PostPermissionListOwnDrafts); err != nil {
		return nil, err
	}

	posts, err := s.repo.GetDraftsByAuthor(ctx, actorUserID)
	if err != nil {
		return nil, normalizeServiceErrorWithOpMsg("post.list_admin_own", "list own draft posts failed", err)
	}
	return posts, nil
}

// GetAdminPostByID returns a single manageable post for the acting user.
func (s *PostService) GetAdminPostByID(ctx context.Context, id uint, actorUserID uint, actorRole string) (entity.Post, error) {
	post, err := s.loadManageablePost(ctx, id, actorUserID, actorRole)
	if err != nil {
		return entity.Post{}, err
	}
	return post, nil
}

// CreateAdminPost persists a new post as Draft.
func (s *PostService) CreateAdminPost(ctx context.Context, actorUserID uint, actorRole string, post entity.Post) (entity.Post, error) {
	if actorUserID == 0 {
		return entity.Post{}, core.ErrPermission
	}
	if err := s.authorizePostAction(ctx, actorRole, core.PostPermissionCreateOwnDraft); err != nil {
		return entity.Post{}, err
	}

	post.AuthorID = actorUserID
	post.Status = entity.StatusDraft

	if err := post.CheckValidity(); err != nil {
		return entity.Post{}, fmt.Errorf("%w: invalid post payload: %v", core.ErrInvalidInput, err)
	}

	generatedSlug := slug.Make(post.Title)
	if generatedSlug == "" {
		return entity.Post{}, fmt.Errorf("%w: title cannot generate a valid slug", core.ErrInvalidInput)
	}

	finalSlug, err := s.generateUniqueSlug(ctx, generatedSlug)
	if err != nil {
		return entity.Post{}, err
	}

	post.Slug = finalSlug

	created, err := s.repo.Create(ctx, post)
	if err != nil {
		return entity.Post{}, normalizeServiceErrorWithOpMsg("post.create_admin", "create admin draft post failed", err)
	}

	if s.media != nil {
		syncCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := s.media.SyncPostReferences(syncCtx, created.ID, created.Content, created.Cover); err != nil {
			log.Printf("[ERROR] Post created (ID: %d) but failed to sync media references: %v", created.ID, err)
		}
	}
	return created, nil
}

func (s *PostService) generateUniqueSlug(ctx context.Context, initialSlug string) (string, error) {
	currentSlug := initialSlug
	counter := 0
	maxAttempts := 100

	for {
		exists, err := s.repo.IsSlugExists(ctx, currentSlug)
		if err != nil {
			return "", normalizeServiceErrorWithOpMsg("post.generate_unique_slug", "check slug uniqueness failed", err)
		}

		if !exists {
			return currentSlug, nil
		}

		counter++
		if counter >= maxAttempts {
			return "", fmt.Errorf("%w: unable to generate unique slug within max attempts", core.ErrConflict)
		}

		currentSlug = fmt.Sprintf("%s-%d", initialSlug, counter)
	}
}

// UpdateAdminPost updates editable post content fields.
func (s *PostService) UpdateAdminPost(ctx context.Context, id uint, patch entity.PostPatch, actorUserID uint, actorRole string) error {
	existingEntity, err := s.loadUpdatablePost(ctx, id, actorUserID, actorRole)
	if err != nil {
		return err
	}

	if patch.Title != nil {
		existingEntity.Title = *patch.Title
	}
	if patch.Content != nil {
		existingEntity.Content = *patch.Content
	}
	if patch.Cover != nil {
		existingEntity.Cover = *patch.Cover
	}
	if patch.CategoryID != nil {
		existingEntity.CategoryID = patch.CategoryID
	}
	if patch.Tags != nil {
		existingEntity.Tags = patch.Tags
	}
	existingEntity.ID = id

	if err := existingEntity.CheckValidity(); err != nil {
		return fmt.Errorf("%w: invalid updated post payload: %v", core.ErrInvalidInput, err)
	}

	if err := s.repo.Update(ctx, existingEntity); err != nil {
		return normalizeServiceErrorWithOpMsg("post.update_admin", "update admin post failed", err)
	}

	if s.media != nil {
		syncCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.media.SyncPostReferences(syncCtx, id, existingEntity.Content, existingEntity.Cover); err != nil {
			log.Printf("[WARN] Post updated (ID: %d) but failed to sync media references: %v", id, err)
		}
	}

	return nil
}

// PublishAdminPost performs the Draft -> Published transition.
func (s *PostService) PublishAdminPost(ctx context.Context, id uint, actorRole string) error {
	if err := s.authorizePostAction(ctx, actorRole, core.PostPermissionPublishPost); err != nil {
		return err
	}

	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return normalizeServiceErrorWithOpMsg("post.publish.load", "load post before publish failed", err)
	}
	if post.Status == entity.StatusPublished {
		return fmt.Errorf("%w: post is already published", core.ErrConflict)
	}
	if err := post.CheckValidity(); err != nil {
		return fmt.Errorf("%w: post is not publishable: %v", core.ErrInvalidInput, err)
	}

	post.Status = entity.StatusPublished
	post.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, post); err != nil {
		return normalizeServiceErrorWithOpMsg("post.publish.update", "persist publish status failed", err)
	}

	return nil
}

// MovePostToDraft performs the minimal "offline" step by moving a post back to Draft.
func (s *PostService) MovePostToDraft(ctx context.Context, id uint, actorRole string) error {
	if err := s.authorizePostAction(ctx, actorRole, core.PostPermissionUnpublishPost); err != nil {
		return err
	}

	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return normalizeServiceErrorWithOpMsg("post.move_draft.load", "load post before move-to-draft failed", err)
	}
	if post.Status == entity.StatusDraft {
		return fmt.Errorf("%w: post is already draft", core.ErrConflict)
	}

	post.Status = entity.StatusDraft
	post.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, post); err != nil {
		return normalizeServiceErrorWithOpMsg("post.move_draft.update", "persist move-to-draft status failed", err)
	}

	return nil
}

// DeleteAdminPost permanently removes a post record.
func (s *PostService) DeleteAdminPost(ctx context.Context, id uint, actorRole string) error {
	if err := s.authorizePostAction(ctx, actorRole, core.PostPermissionDeletePost); err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return normalizeServiceErrorWithOpMsg("post.delete_admin", "delete admin post failed", err)
	}
	return nil
}

func (s *PostService) loadManageablePost(ctx context.Context, id uint, actorUserID uint, actorRole string) (entity.Post, error) {
	canReadAny, err := s.hasPostPermission(ctx, actorRole, core.PostPermissionReadAnyPost)
	if err != nil {
		return entity.Post{}, err
	}
	if canReadAny {
		post, err := s.repo.GetByID(ctx, id)
		if err != nil {
			return entity.Post{}, normalizeServiceErrorWithOpMsg("post.load_manageable", "load manageable post failed", err)
		}
		return post, nil
	}

	if actorUserID == 0 {
		return entity.Post{}, core.ErrPermission
	}
	if err := s.authorizePostAction(ctx, actorRole, core.PostPermissionReadOwnDraft); err != nil {
		return entity.Post{}, err
	}

	post, err := s.repo.GetDraftByIDAndAuthor(ctx, id, actorUserID)
	if err != nil {
		return entity.Post{}, normalizeServiceErrorWithOpMsg("post.load_manageable_own", "load own manageable draft failed", err)
	}
	return post, nil
}

func (s *PostService) loadUpdatablePost(ctx context.Context, id uint, actorUserID uint, actorRole string) (entity.Post, error) {
	canUpdateAny, err := s.hasPostPermission(ctx, actorRole, core.PostPermissionUpdateAnyPost)
	if err != nil {
		return entity.Post{}, err
	}
	if canUpdateAny {
		post, err := s.repo.GetByID(ctx, id)
		if err != nil {
			return entity.Post{}, normalizeServiceErrorWithOpMsg("post.load_updatable", "load updatable post failed", err)
		}
		return post, nil
	}

	if actorUserID == 0 {
		return entity.Post{}, core.ErrPermission
	}
	if err := s.authorizePostAction(ctx, actorRole, core.PostPermissionUpdateOwnDraft); err != nil {
		return entity.Post{}, err
	}

	post, err := s.repo.GetDraftByIDAndAuthor(ctx, id, actorUserID)
	if err != nil {
		return entity.Post{}, normalizeServiceErrorWithOpMsg("post.load_updatable_own", "load own updatable draft failed", err)
	}
	return post, nil
}

func (s *PostService) authorizePostAction(ctx context.Context, actorRole string, permission core.PostPermission) error {
	if s.authorizer == nil {
		return core.ErrPermission
	}
	if err := s.authorizer.AuthorizePostAction(ctx, actorRole, permission); err != nil {
		return normalizeServiceErrorWithOpMsg("post.authorize", "authorize post action failed", err)
	}
	return nil
}

func (s *PostService) hasPostPermission(ctx context.Context, actorRole string, permission core.PostPermission) (bool, error) {
	if err := s.authorizePostAction(ctx, actorRole, permission); err != nil {
		if errors.Is(err, core.ErrPermission) {
			return false, nil
		}
		return false, normalizeServiceErrorWithOpMsg("post.has_permission", "check post permission failed", err)
	}
	return true, nil
}
