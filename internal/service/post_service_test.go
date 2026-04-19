package service

import (
	"context"
	"errors"
	"testing"

	"KaldalisCMS/internal/core"
	"KaldalisCMS/internal/core/entity"
)

type fakePostRepo struct {
	getByIDFn              func(ctx context.Context, id uint) (entity.Post, error)
	getPublishedByIDFn     func(ctx context.Context, id uint) (entity.Post, error)
	getDraftByIDAndAuthorFn func(ctx context.Context, id uint, authorID uint) (entity.Post, error)
	createFn               func(ctx context.Context, post entity.Post) (entity.Post, error)
	updateFn               func(ctx context.Context, post entity.Post) error
	deleteFn               func(ctx context.Context, id uint) error
	getAllFn               func(ctx context.Context) ([]entity.Post, error)
	getPublishedFn         func(ctx context.Context) ([]entity.Post, error)
	getDraftsByAuthorFn    func(ctx context.Context, authorID uint) ([]entity.Post, error)
	isSlugExistsFn         func(ctx context.Context, slug string) (bool, error)
}

func (f *fakePostRepo) GetByID(ctx context.Context, id uint) (entity.Post, error) {
	return f.getByIDFn(ctx, id)
}
func (f *fakePostRepo) GetPublishedByID(ctx context.Context, id uint) (entity.Post, error) {
	return f.getPublishedByIDFn(ctx, id)
}
func (f *fakePostRepo) GetDraftByIDAndAuthor(ctx context.Context, id uint, authorID uint) (entity.Post, error) {
	return f.getDraftByIDAndAuthorFn(ctx, id, authorID)
}
func (f *fakePostRepo) Create(ctx context.Context, post entity.Post) (entity.Post, error) {
	return f.createFn(ctx, post)
}
func (f *fakePostRepo) Update(ctx context.Context, post entity.Post) error {
	return f.updateFn(ctx, post)
}
func (f *fakePostRepo) Delete(ctx context.Context, id uint) error { return f.deleteFn(ctx, id) }
func (f *fakePostRepo) GetAll(ctx context.Context) ([]entity.Post, error) {
	return f.getAllFn(ctx)
}
func (f *fakePostRepo) GetPublished(ctx context.Context) ([]entity.Post, error) {
	return f.getPublishedFn(ctx)
}
func (f *fakePostRepo) GetDraftsByAuthor(ctx context.Context, authorID uint) ([]entity.Post, error) {
	return f.getDraftsByAuthorFn(ctx, authorID)
}
func (f *fakePostRepo) IsSlugExists(ctx context.Context, slug string) (bool, error) {
	return f.isSlugExistsFn(ctx, slug)
}

// fakeAuthorizer grants the exact permissions in `allow`. Others return ErrPermission.
type fakeAuthorizer struct {
	allow map[core.PostPermission]bool
	err   error // if set, returned instead of checking allow
}

func (f *fakeAuthorizer) AuthorizePostAction(ctx context.Context, role string, p core.PostPermission) error {
	if f.err != nil {
		return f.err
	}
	if f.allow[p] {
		return nil
	}
	return core.ErrPermission
}

func allowAll() *fakeAuthorizer {
	return &fakeAuthorizer{allow: map[core.PostPermission]bool{
		core.PostPermissionCreateOwnDraft: true,
		core.PostPermissionListOwnDrafts:  true,
		core.PostPermissionReadOwnDraft:   true,
		core.PostPermissionUpdateOwnDraft: true,
		core.PostPermissionListAnyPost:    true,
		core.PostPermissionReadAnyPost:    true,
		core.PostPermissionUpdateAnyPost:  true,
		core.PostPermissionPublishPost:    true,
		core.PostPermissionUnpublishPost:  true,
		core.PostPermissionDeletePost:     true,
	}}
}

func TestPostService_ListPublicPosts(t *testing.T) {
	ctx := context.Background()
	repo := &fakePostRepo{getPublishedFn: func(ctx context.Context) ([]entity.Post, error) {
		return []entity.Post{{ID: 1}}, nil
	}}
	got, err := NewPostService(repo, allowAll()).ListPublicPosts(ctx)
	if err != nil || len(got) != 1 {
		t.Fatalf("unexpected: %+v %v", got, err)
	}
}

func TestPostService_GetPublicPostByID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &fakePostRepo{getPublishedByIDFn: func(ctx context.Context, id uint) (entity.Post, error) {
		return entity.Post{}, core.ErrNotFound
	}}
	_, err := NewPostService(repo, allowAll()).GetPublicPostByID(ctx, 1)
	if !errors.Is(err, core.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestPostService_ListAdminPosts_AdminScope(t *testing.T) {
	ctx := context.Background()
	repo := &fakePostRepo{getAllFn: func(ctx context.Context) ([]entity.Post, error) {
		return []entity.Post{{ID: 1}, {ID: 2}}, nil
	}}
	got, err := NewPostService(repo, allowAll()).ListAdminPosts(ctx, 9, "admin")
	if err != nil || len(got) != 2 {
		t.Fatalf("unexpected: %+v %v", got, err)
	}
}

func TestPostService_ListAdminPosts_OwnDraftsScope(t *testing.T) {
	ctx := context.Background()
	auth := &fakeAuthorizer{allow: map[core.PostPermission]bool{
		core.PostPermissionListOwnDrafts: true,
	}}
	repo := &fakePostRepo{getDraftsByAuthorFn: func(ctx context.Context, authorID uint) ([]entity.Post, error) {
		if authorID != 5 {
			t.Fatalf("want authorID 5, got %d", authorID)
		}
		return []entity.Post{{ID: 10}}, nil
	}}
	got, err := NewPostService(repo, auth).ListAdminPosts(ctx, 5, "editor")
	if err != nil || len(got) != 1 {
		t.Fatalf("unexpected: %+v %v", got, err)
	}
}

func TestPostService_ListAdminPosts_AnonymousRejected(t *testing.T) {
	ctx := context.Background()
	auth := &fakeAuthorizer{allow: map[core.PostPermission]bool{}}
	_, err := NewPostService(&fakePostRepo{}, auth).ListAdminPosts(ctx, 0, "guest")
	if !errors.Is(err, core.ErrPermission) {
		t.Fatalf("want ErrPermission, got %v", err)
	}
}

func TestPostService_CreateAdminPost_HappyPath(t *testing.T) {
	ctx := context.Background()
	slugExists := map[string]bool{"hello-world": true} // force one retry
	var created entity.Post
	repo := &fakePostRepo{
		isSlugExistsFn: func(ctx context.Context, slug string) (bool, error) {
			return slugExists[slug], nil
		},
		createFn: func(ctx context.Context, p entity.Post) (entity.Post, error) {
			p.ID = 100
			created = p
			return p, nil
		},
	}
	got, err := NewPostService(repo, allowAll()).CreateAdminPost(ctx, 7, "admin", entity.Post{
		Title:   "Hello World",
		Content: "body",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 100 {
		t.Fatalf("ID not returned: %+v", got)
	}
	if created.AuthorID != 7 {
		t.Fatalf("author not set: %d", created.AuthorID)
	}
	if created.Status != entity.StatusDraft {
		t.Fatalf("status should be draft, got %d", created.Status)
	}
	if created.Slug != "hello-world-1" {
		t.Fatalf("slug retry failed: %q", created.Slug)
	}
}

func TestPostService_CreateAdminPost_EmptyTitle(t *testing.T) {
	ctx := context.Background()
	_, err := NewPostService(&fakePostRepo{}, allowAll()).CreateAdminPost(ctx, 1, "admin", entity.Post{})
	if !errors.Is(err, core.ErrInvalidInput) {
		t.Fatalf("want ErrInvalidInput, got %v", err)
	}
}

func TestPostService_CreateAdminPost_Anonymous(t *testing.T) {
	ctx := context.Background()
	_, err := NewPostService(&fakePostRepo{}, allowAll()).CreateAdminPost(ctx, 0, "admin", entity.Post{Title: "x"})
	if !errors.Is(err, core.ErrPermission) {
		t.Fatalf("want ErrPermission, got %v", err)
	}
}

func TestPostService_CreateAdminPost_NoPermission(t *testing.T) {
	ctx := context.Background()
	auth := &fakeAuthorizer{allow: map[core.PostPermission]bool{}}
	_, err := NewPostService(&fakePostRepo{}, auth).CreateAdminPost(ctx, 1, "guest", entity.Post{Title: "x"})
	if !errors.Is(err, core.ErrPermission) {
		t.Fatalf("want ErrPermission, got %v", err)
	}
}

func TestPostService_PublishAdminPost_AlreadyPublished(t *testing.T) {
	ctx := context.Background()
	repo := &fakePostRepo{getByIDFn: func(ctx context.Context, id uint) (entity.Post, error) {
		return entity.Post{ID: id, Title: "t", Status: entity.StatusPublished}, nil
	}}
	err := NewPostService(repo, allowAll()).PublishAdminPost(ctx, 1, 9, "admin")
	if !errors.Is(err, core.ErrConflict) {
		t.Fatalf("want ErrConflict, got %v", err)
	}
}

func TestPostService_PublishAdminPost_Success(t *testing.T) {
	ctx := context.Background()
	var updated entity.Post
	repo := &fakePostRepo{
		getByIDFn: func(ctx context.Context, id uint) (entity.Post, error) {
			return entity.Post{ID: id, Title: "t", Status: entity.StatusDraft}, nil
		},
		updateFn: func(ctx context.Context, p entity.Post) error {
			updated = p
			return nil
		},
	}
	if err := NewPostService(repo, allowAll()).PublishAdminPost(ctx, 1, 9, "admin"); err != nil {
		t.Fatal(err)
	}
	if updated.Status != entity.StatusPublished {
		t.Fatalf("status not published: %d", updated.Status)
	}
}

func TestPostService_PublishAdminPost_NoPermission(t *testing.T) {
	ctx := context.Background()
	auth := &fakeAuthorizer{allow: map[core.PostPermission]bool{}}
	err := NewPostService(&fakePostRepo{}, auth).PublishAdminPost(ctx, 1, 9, "guest")
	if !errors.Is(err, core.ErrPermission) {
		t.Fatalf("want ErrPermission, got %v", err)
	}
}

func TestPostService_MovePostToDraft_AlreadyDraft(t *testing.T) {
	ctx := context.Background()
	repo := &fakePostRepo{getByIDFn: func(ctx context.Context, id uint) (entity.Post, error) {
		return entity.Post{ID: id, Title: "t", Status: entity.StatusDraft}, nil
	}}
	err := NewPostService(repo, allowAll()).MovePostToDraft(ctx, 1, 9, "admin")
	if !errors.Is(err, core.ErrConflict) {
		t.Fatalf("want ErrConflict, got %v", err)
	}
}

func TestPostService_DeleteAdminPost_Success(t *testing.T) {
	ctx := context.Background()
	var deletedID uint
	repo := &fakePostRepo{
		getByIDFn: func(ctx context.Context, id uint) (entity.Post, error) {
			return entity.Post{ID: id, Title: "t"}, nil
		},
		deleteFn: func(ctx context.Context, id uint) error {
			deletedID = id
			return nil
		},
	}
	if err := NewPostService(repo, allowAll()).DeleteAdminPost(ctx, 42, 9, "admin"); err != nil {
		t.Fatal(err)
	}
	if deletedID != 42 {
		t.Fatalf("want 42, got %d", deletedID)
	}
}

func TestPostService_UpdateAdminPost_PatchApplied(t *testing.T) {
	ctx := context.Background()
	var updated entity.Post
	repo := &fakePostRepo{
		getByIDFn: func(ctx context.Context, id uint) (entity.Post, error) {
			return entity.Post{ID: id, Title: "old", Content: "old body"}, nil
		},
		updateFn: func(ctx context.Context, p entity.Post) error {
			updated = p
			return nil
		},
	}
	newTitle := "new title"
	newContent := "new body"
	err := NewPostService(repo, allowAll()).UpdateAdminPost(ctx, 1, entity.PostPatch{
		Title:   &newTitle,
		Content: &newContent,
	}, 9, "admin")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != newTitle || updated.Content != newContent {
		t.Fatalf("patch not applied: %+v", updated)
	}
}

func TestPostService_GetAdminPostByID_OwnDraftScope(t *testing.T) {
	ctx := context.Background()
	auth := &fakeAuthorizer{allow: map[core.PostPermission]bool{
		core.PostPermissionReadOwnDraft: true,
	}}
	repo := &fakePostRepo{
		getDraftByIDAndAuthorFn: func(ctx context.Context, id uint, authorID uint) (entity.Post, error) {
			if authorID != 5 {
				t.Fatalf("wrong authorID %d", authorID)
			}
			return entity.Post{ID: id, AuthorID: authorID}, nil
		},
	}
	got, err := NewPostService(repo, auth).GetAdminPostByID(ctx, 3, 5, "editor")
	if err != nil || got.ID != 3 {
		t.Fatalf("unexpected: %+v %v", got, err)
	}
}
