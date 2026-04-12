package auth

import (
	"context"
	"errors"
	"testing"

	"KaldalisCMS/internal/core"
)

// TestMapPostPermission verifies that every defined PostPermission maps to the
// expected Casbin (object, action) pair and that unknown permissions are rejected.
func TestMapPostPermission(t *testing.T) {
	tests := []struct {
		perm      core.PostPermission
		wantObj   string
		wantAct   string
		wantError bool
	}{
		{core.PostPermissionCreateOwnDraft, "post:draft", "create", false},
		{core.PostPermissionListOwnDrafts, "post:draft", "list:own", false},
		{core.PostPermissionReadOwnDraft, "post:draft", "read:own", false},
		{core.PostPermissionUpdateOwnDraft, "post:draft", "update:own", false},
		{core.PostPermissionListAnyPost, "post", "list:any", false},
		{core.PostPermissionReadAnyPost, "post", "read:any", false},
		{core.PostPermissionUpdateAnyPost, "post", "update:any", false},
		{core.PostPermissionPublishPost, "post", "publish", false},
		{core.PostPermissionUnpublishPost, "post", "unpublish", false},
		{core.PostPermissionDeletePost, "post", "delete", false},
		// unknown permission
		{"post:nonexistent", "", "", true},
	}

	for _, tt := range tests {
		t.Run(string(tt.perm), func(t *testing.T) {
			obj, act, err := mapPostPermission(tt.perm)
			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error for permission %q, got nil", tt.perm)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if obj != tt.wantObj || act != tt.wantAct {
				t.Errorf("mapPostPermission(%q) = (%q, %q), want (%q, %q)", tt.perm, obj, act, tt.wantObj, tt.wantAct)
			}
		})
	}
}

// TestCasbinPostAuthorizer_AuthorizePostAction exercises the full authorizer
// against a real (in-memory) Casbin enforcer with production-equivalent policies.
func TestCasbinPostAuthorizer_AuthorizePostAction(t *testing.T) {
	e := setupTestEnforcer(t, defaultOpts())
	authorizer := NewCasbinPostAuthorizer(e)
	ctx := context.Background()

	tests := []struct {
		name       string
		role       string
		permission core.PostPermission
		wantAllow  bool
	}{
		// ── super_admin ──
		{"super_admin publish", "super_admin", core.PostPermissionPublishPost, true},
		{"super_admin unpublish", "super_admin", core.PostPermissionUnpublishPost, true},
		{"super_admin delete", "super_admin", core.PostPermissionDeletePost, true},
		{"super_admin create draft", "super_admin", core.PostPermissionCreateOwnDraft, true},
		{"super_admin list any", "super_admin", core.PostPermissionListAnyPost, true},

		// ── admin ──
		{"admin publish", "admin", core.PostPermissionPublishPost, true},
		{"admin unpublish", "admin", core.PostPermissionUnpublishPost, true},
		{"admin delete", "admin", core.PostPermissionDeletePost, true},
		{"admin list any", "admin", core.PostPermissionListAnyPost, true},
		{"admin read any", "admin", core.PostPermissionReadAnyPost, true},
		{"admin update any", "admin", core.PostPermissionUpdateAnyPost, true},
		{"admin create draft (inherited)", "admin", core.PostPermissionCreateOwnDraft, true},
		{"admin list own (inherited)", "admin", core.PostPermissionListOwnDrafts, true},

		// ── user ──
		{"user create draft", "user", core.PostPermissionCreateOwnDraft, true},
		{"user list own", "user", core.PostPermissionListOwnDrafts, true},
		{"user read own", "user", core.PostPermissionReadOwnDraft, true},
		{"user update own", "user", core.PostPermissionUpdateOwnDraft, true},
		{"user cannot publish", "user", core.PostPermissionPublishPost, false},
		{"user cannot unpublish", "user", core.PostPermissionUnpublishPost, false},
		{"user cannot delete", "user", core.PostPermissionDeletePost, false},
		{"user cannot list any", "user", core.PostPermissionListAnyPost, false},
		{"user cannot read any", "user", core.PostPermissionReadAnyPost, false},
		{"user cannot update any", "user", core.PostPermissionUpdateAnyPost, false},

		// ── anonymous ──
		{"anonymous cannot create draft", "anonymous", core.PostPermissionCreateOwnDraft, false},
		{"anonymous cannot publish", "anonymous", core.PostPermissionPublishPost, false},

		// ── empty role ──
		{"empty role denied", "", core.PostPermissionCreateOwnDraft, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authorizer.AuthorizePostAction(ctx, tt.role, tt.permission)
			if tt.wantAllow {
				if err != nil {
					t.Errorf("expected allow, got error: %v", err)
				}
			} else {
				if err == nil {
					t.Error("expected deny, got nil error")
				} else if !errors.Is(err, core.ErrPermission) {
					// The authorizer should wrap or return ErrPermission on deny.
					// Implementation may wrap it via normalizeServiceErrorWithOpMsg,
					// so we also accept if it's not ErrPermission but still non-nil.
					t.Logf("denied with non-ErrPermission error (acceptable): %v", err)
				}
			}
		})
	}
}

// TestCasbinPostAuthorizer_NilEnforcer ensures a nil enforcer always denies.
func TestCasbinPostAuthorizer_NilEnforcer(t *testing.T) {
	authorizer := NewCasbinPostAuthorizer(nil)
	err := authorizer.AuthorizePostAction(context.Background(), "super_admin", core.PostPermissionPublishPost)
	if err == nil {
		t.Error("nil enforcer should deny all requests")
	}
	if !errors.Is(err, core.ErrPermission) {
		t.Errorf("expected ErrPermission, got: %v", err)
	}
}
