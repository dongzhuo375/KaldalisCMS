package v1

import (
	"context"

	"KaldalisCMS/internal/core/entity"
)

// fakePostService implements core.PostService for handler-layer tests.
// Each field is an optional override; nil fields cause the test to panic if hit,
// which surfaces accidentally-exercised branches instead of silently passing.
type fakePostService struct {
	listPublicFn       func(ctx context.Context) ([]entity.Post, error)
	getPublicByIDFn    func(ctx context.Context, id uint) (entity.Post, error)
	listAdminFn        func(ctx context.Context, uid uint, role string) ([]entity.Post, error)
	getAdminByIDFn     func(ctx context.Context, id uint, uid uint, role string) (entity.Post, error)
	createAdminFn      func(ctx context.Context, uid uint, role string, p entity.Post) (entity.Post, error)
	updateAdminFn      func(ctx context.Context, id uint, patch entity.PostPatch, uid uint, role string) error
	deleteAdminFn      func(ctx context.Context, id uint, uid uint, role string) error
	publishAdminFn     func(ctx context.Context, id uint, uid uint, role string) error
	moveToDraftAdminFn func(ctx context.Context, id uint, uid uint, role string) error
}

func (f *fakePostService) ListPublicPosts(ctx context.Context) ([]entity.Post, error) {
	return f.listPublicFn(ctx)
}
func (f *fakePostService) GetPublicPostByID(ctx context.Context, id uint) (entity.Post, error) {
	return f.getPublicByIDFn(ctx, id)
}
func (f *fakePostService) ListAdminPosts(ctx context.Context, uid uint, role string) ([]entity.Post, error) {
	return f.listAdminFn(ctx, uid, role)
}
func (f *fakePostService) GetAdminPostByID(ctx context.Context, id uint, uid uint, role string) (entity.Post, error) {
	return f.getAdminByIDFn(ctx, id, uid, role)
}
func (f *fakePostService) CreateAdminPost(ctx context.Context, uid uint, role string, p entity.Post) (entity.Post, error) {
	return f.createAdminFn(ctx, uid, role, p)
}
func (f *fakePostService) UpdateAdminPost(ctx context.Context, id uint, patch entity.PostPatch, uid uint, role string) error {
	return f.updateAdminFn(ctx, id, patch, uid, role)
}
func (f *fakePostService) DeleteAdminPost(ctx context.Context, id uint, uid uint, role string) error {
	return f.deleteAdminFn(ctx, id, uid, role)
}
func (f *fakePostService) PublishAdminPost(ctx context.Context, id uint, uid uint, role string) error {
	return f.publishAdminFn(ctx, id, uid, role)
}
func (f *fakePostService) MovePostToDraft(ctx context.Context, id uint, uid uint, role string) error {
	return f.moveToDraftAdminFn(ctx, id, uid, role)
}
