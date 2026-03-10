package auth

import (
	"KaldalisCMS/internal/core"
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
)

// CasbinPostAuthorizer adapts post workflow capabilities onto Casbin policies.
// It keeps object/action mapping in the infrastructure layer so services can speak in domain terms.
type CasbinPostAuthorizer struct {
	enforcer *casbin.Enforcer
}

func NewCasbinPostAuthorizer(enforcer *casbin.Enforcer) *CasbinPostAuthorizer {
	return &CasbinPostAuthorizer{enforcer: enforcer}
}

func (a *CasbinPostAuthorizer) AuthorizePostAction(ctx context.Context, role string, permission core.PostPermission) error {
	_, _ = ctx, permission
	if a.enforcer == nil || role == "" {
		return core.ErrPermission
	}

	obj, act, err := mapPostPermission(permission)
	if err != nil {
		return err
	}

	ok, err := a.enforcer.Enforce(role, obj, act)
	if err != nil {
		return fmt.Errorf("post authorization failed: %w", err)
	}
	if !ok {
		return core.ErrPermission
	}
	return nil
}

func mapPostPermission(permission core.PostPermission) (obj string, act string, err error) {
	switch permission {
	case core.PostPermissionCreateOwnDraft:
		return "post:draft", "create", nil
	case core.PostPermissionListOwnDrafts:
		return "post:draft", "list:own", nil
	case core.PostPermissionReadOwnDraft:
		return "post:draft", "read:own", nil
	case core.PostPermissionUpdateOwnDraft:
		return "post:draft", "update:own", nil
	case core.PostPermissionListAnyPost:
		return "post", "list:any", nil
	case core.PostPermissionReadAnyPost:
		return "post", "read:any", nil
	case core.PostPermissionUpdateAnyPost:
		return "post", "update:any", nil
	case core.PostPermissionPublishPost:
		return "post", "publish", nil
	case core.PostPermissionUnpublishPost:
		return "post", "unpublish", nil
	case core.PostPermissionDeletePost:
		return "post", "delete", nil
	default:
		return "", "", fmt.Errorf("unknown post permission: %s", permission)
	}
}
