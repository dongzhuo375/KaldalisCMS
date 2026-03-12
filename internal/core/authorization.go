package core

import "context"

// PostPermission represents a post-management capability understood by the authorization layer.
// Service code depends on these business-level permissions instead of hard-coded role names.
type PostPermission string

const (
	PostPermissionCreateOwnDraft PostPermission = "post:create_own_draft"
	PostPermissionListOwnDrafts  PostPermission = "post:list_own_drafts"
	PostPermissionReadOwnDraft   PostPermission = "post:read_own_draft"
	PostPermissionUpdateOwnDraft PostPermission = "post:update_own_draft"
	PostPermissionListAnyPost    PostPermission = "post:list_any"
	PostPermissionReadAnyPost    PostPermission = "post:read_any"
	PostPermissionUpdateAnyPost  PostPermission = "post:update_any"
	PostPermissionPublishPost    PostPermission = "post:publish"
	PostPermissionUnpublishPost  PostPermission = "post:unpublish"
	PostPermissionDeletePost     PostPermission = "post:delete"
)

// PostAuthorizer decides whether a role currently has a given post-management capability.
// Infrastructure implementations may use Casbin, another policy engine, or a remote authorization service.
type PostAuthorizer interface {
	AuthorizePostAction(ctx context.Context, role string, permission PostPermission) error
}
