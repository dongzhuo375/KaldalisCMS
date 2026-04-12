package auth

import (
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

// modelText is the production Casbin model (casbin_model.conf) inlined for test isolation.
const modelText = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == "super_admin" || (g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (p.act == "*" || r.act == p.act))
`

// setupTestEnforcer creates a Casbin enforcer with the production model and
// the same policy set that setup_service.go seeds during installation.
// The flags mirror the setup wizard checkboxes.
type setupOpts struct {
	AdminCanDelete     bool
	UserCanUpload      bool
	AllowAnonymousRead bool
}

func defaultOpts() setupOpts {
	return setupOpts{
		AdminCanDelete:     true,
		UserCanUpload:      false,
		AllowAnonymousRead: true,
	}
}

func setupTestEnforcer(t *testing.T, opts setupOpts) *casbin.Enforcer {
	t.Helper()

	m, err := model.NewModelFromString(modelText)
	if err != nil {
		t.Fatalf("failed to load model: %v", err)
	}

	// Use an empty file adapter (no persistent file); policies are added programmatically.
	e, err := casbin.NewEnforcer(m, fileadapter.NewAdapter(""))
	if err != nil {
		t.Fatalf("failed to create enforcer: %v", err)
	}

	// ── Policies identical to setup_service.go Install() ──
	// super_admin has no explicit policies; the matcher hardcodes the bypass.

	// 1. admin — route policies
	adminRoutes := [][]string{
		{"admin", "/api/v1/admin/posts", "GET"},
		{"admin", "/api/v1/admin/posts", "POST"},
		{"admin", "/api/v1/admin/posts/:id", "GET"},
		{"admin", "/api/v1/admin/posts/:id", "PUT"},
		{"admin", "/api/v1/admin/posts/:id/publish", "POST"},
		{"admin", "/api/v1/admin/posts/:id/draft", "POST"},
		// capability policies
		{"admin", "post", "list:any"},
		{"admin", "post", "read:any"},
		{"admin", "post", "update:any"},
		{"admin", "post", "publish"},
		{"admin", "post", "unpublish"},
		{"admin", "post", "delete"},
		// media / tags / categories
		{"admin", "/api/v1/media", "POST"},
		{"admin", "/api/v1/tags", "POST"},
		{"admin", "/api/v1/tags/:id", "PUT"},
		{"admin", "/api/v1/categories", "POST"},
		{"admin", "/api/v1/categories/:id", "PUT"},
	}
	_, _ = e.AddPolicies(adminRoutes)

	if opts.AdminCanDelete {
		_, _ = e.AddPolicy("admin", "/api/v1/admin/posts/:id", "DELETE")
		_, _ = e.AddPolicy("admin", "/api/v1/media/:id", "DELETE")
		_, _ = e.AddPolicy("admin", "/api/v1/tags/:id", "DELETE")
		_, _ = e.AddPolicy("admin", "/api/v1/categories/:id", "DELETE")
	}

	// 3. user — route policies
	userRoutes := [][]string{
		{"user", "/api/v1/posts", "GET"},
		{"user", "/api/v1/posts/:id", "GET"},
		{"user", "/api/v1/admin/posts", "GET"},
		{"user", "/api/v1/admin/posts", "POST"},
		{"user", "/api/v1/admin/posts/:id", "GET"},
		{"user", "/api/v1/admin/posts/:id", "PUT"},
		// capability policies
		{"user", "post:draft", "create"},
		{"user", "post:draft", "list:own"},
		{"user", "post:draft", "read:own"},
		{"user", "post:draft", "update:own"},
		// media read
		{"user", "/api/v1/media", "GET"},
	}
	_, _ = e.AddPolicies(userRoutes)

	if opts.UserCanUpload {
		_, _ = e.AddPolicy("user", "/api/v1/media", "POST")
	}

	// 4. anonymous
	if opts.AllowAnonymousRead {
		_, _ = e.AddPolicy("anonymous", "/api/v1/posts", "GET")
		_, _ = e.AddPolicy("anonymous", "/api/v1/posts/:id", "GET")
	}

	// 5. Role inheritance
	_, _ = e.AddGroupingPolicy("admin", "user")
	_, _ = e.AddGroupingPolicy("super_admin", "admin")

	// 6. Logout policies (from ensurePostWorkflowPolicies)
	_, _ = e.AddPolicy("admin", "/api/v1/users/logout", "POST")
	_, _ = e.AddPolicy("user", "/api/v1/users/logout", "POST")
	_, _ = e.AddPolicy("super_admin", "/api/v1/users/logout", "POST")

	return e
}

func enforce(t *testing.T, e *casbin.Enforcer, sub, obj, act string) bool {
	t.Helper()
	ok, err := e.Enforce(sub, obj, act)
	if err != nil {
		t.Fatalf("Enforce(%q, %q, %q) error: %v", sub, obj, act, err)
	}
	return ok
}

// ═══════════════════════════════════════════════════════════════════════
// Test: Route-level policy matrix (what the HTTP middleware checks)
// ═══════════════════════════════════════════════════════════════════════

func TestRoutePolicy_DefaultSetup(t *testing.T) {
	e := setupTestEnforcer(t, defaultOpts())

	tests := []struct {
		name   string
		sub    string
		obj    string
		act    string
		expect bool
	}{
		// ── super_admin: bypass all ──
		{"super_admin can GET admin posts", "super_admin", "/api/v1/admin/posts", "GET", true},
		{"super_admin can POST admin posts", "super_admin", "/api/v1/admin/posts", "POST", true},
		{"super_admin can DELETE admin post", "super_admin", "/api/v1/admin/posts/:id", "DELETE", true},
		{"super_admin can publish", "super_admin", "/api/v1/admin/posts/:id/publish", "POST", true},
		{"super_admin can access any path", "super_admin", "/api/v1/anything/here", "PATCH", true},

		// ── admin: explicit route policies ──
		{"admin can GET admin posts", "admin", "/api/v1/admin/posts", "GET", true},
		{"admin can POST admin posts", "admin", "/api/v1/admin/posts", "POST", true},
		{"admin can GET admin post by id", "admin", "/api/v1/admin/posts/:id", "GET", true},
		{"admin can PUT admin post", "admin", "/api/v1/admin/posts/:id", "PUT", true},
		{"admin can DELETE admin post", "admin", "/api/v1/admin/posts/:id", "DELETE", true},
		{"admin can publish post", "admin", "/api/v1/admin/posts/:id/publish", "POST", true},
		{"admin can draft post", "admin", "/api/v1/admin/posts/:id/draft", "POST", true},
		{"admin can POST media", "admin", "/api/v1/media", "POST", true},
		{"admin can DELETE media", "admin", "/api/v1/media/:id", "DELETE", true},
		{"admin can logout", "admin", "/api/v1/users/logout", "POST", true},
		// admin inherits user's public read
		{"admin inherits user GET posts", "admin", "/api/v1/posts", "GET", true},
		{"admin inherits user GET post by id", "admin", "/api/v1/posts/:id", "GET", true},
		{"admin inherits user GET media", "admin", "/api/v1/media", "GET", true},
		// admin cannot PATCH (no such policy)
		{"admin cannot PATCH admin posts", "admin", "/api/v1/admin/posts/:id", "PATCH", false},

		// ── user: limited access ──
		{"user can GET public posts", "user", "/api/v1/posts", "GET", true},
		{"user can GET public post by id", "user", "/api/v1/posts/:id", "GET", true},
		{"user can GET admin posts (own drafts)", "user", "/api/v1/admin/posts", "GET", true},
		{"user can POST admin posts (create draft)", "user", "/api/v1/admin/posts", "POST", true},
		{"user can GET admin post by id", "user", "/api/v1/admin/posts/:id", "GET", true},
		{"user can PUT admin post (own draft)", "user", "/api/v1/admin/posts/:id", "PUT", true},
		{"user can GET media", "user", "/api/v1/media", "GET", true},
		{"user can logout", "user", "/api/v1/users/logout", "POST", true},
		// user CANNOT access publish/draft/delete routes
		{"user cannot publish post", "user", "/api/v1/admin/posts/:id/publish", "POST", false},
		{"user cannot draft post", "user", "/api/v1/admin/posts/:id/draft", "POST", false},
		{"user cannot DELETE admin post", "user", "/api/v1/admin/posts/:id", "DELETE", false},
		{"user cannot POST media (no upload)", "user", "/api/v1/media", "POST", false},
		{"user cannot DELETE media", "user", "/api/v1/media/:id", "DELETE", false},

		// ── anonymous: only public read ──
		{"anonymous can GET public posts", "anonymous", "/api/v1/posts", "GET", true},
		{"anonymous can GET public post by id", "anonymous", "/api/v1/posts/:id", "GET", true},
		{"anonymous cannot GET admin posts", "anonymous", "/api/v1/admin/posts", "GET", false},
		{"anonymous cannot POST admin posts", "anonymous", "/api/v1/admin/posts", "POST", false},
		{"anonymous cannot DELETE", "anonymous", "/api/v1/admin/posts/:id", "DELETE", false},
		{"anonymous cannot publish", "anonymous", "/api/v1/admin/posts/:id/publish", "POST", false},
		{"anonymous cannot logout", "anonymous", "/api/v1/users/logout", "POST", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := enforce(t, e, tt.sub, tt.obj, tt.act)
			if got != tt.expect {
				t.Errorf("Enforce(%q, %q, %q) = %v, want %v", tt.sub, tt.obj, tt.act, got, tt.expect)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════════════
// Test: Capability policy matrix (what the service-layer PostAuthorizer checks)
// ═══════════════════════════════════════════════════════════════════════

func TestCapabilityPolicy_DefaultSetup(t *testing.T) {
	e := setupTestEnforcer(t, defaultOpts())

	tests := []struct {
		name   string
		sub    string
		obj    string
		act    string
		expect bool
	}{
		// ── super_admin: bypass ──
		{"super_admin can publish", "super_admin", "post", "publish", true},
		{"super_admin can unpublish", "super_admin", "post", "unpublish", true},
		{"super_admin can delete", "super_admin", "post", "delete", true},
		{"super_admin can create draft", "super_admin", "post:draft", "create", true},

		// ── admin: full post management ──
		{"admin can list:any", "admin", "post", "list:any", true},
		{"admin can read:any", "admin", "post", "read:any", true},
		{"admin can update:any", "admin", "post", "update:any", true},
		{"admin can publish", "admin", "post", "publish", true},
		{"admin can unpublish", "admin", "post", "unpublish", true},
		{"admin can delete", "admin", "post", "delete", true},
		// admin inherits user's draft capabilities
		{"admin can create draft (inherited)", "admin", "post:draft", "create", true},
		{"admin can list:own (inherited)", "admin", "post:draft", "list:own", true},
		{"admin can read:own (inherited)", "admin", "post:draft", "read:own", true},
		{"admin can update:own (inherited)", "admin", "post:draft", "update:own", true},

		// ── user: only own drafts ──
		{"user can create draft", "user", "post:draft", "create", true},
		{"user can list:own", "user", "post:draft", "list:own", true},
		{"user can read:own", "user", "post:draft", "read:own", true},
		{"user can update:own", "user", "post:draft", "update:own", true},
		// user CANNOT do admin-level operations
		{"user cannot list:any", "user", "post", "list:any", false},
		{"user cannot read:any", "user", "post", "read:any", false},
		{"user cannot update:any", "user", "post", "update:any", false},
		{"user cannot publish", "user", "post", "publish", false},
		{"user cannot unpublish", "user", "post", "unpublish", false},
		{"user cannot delete", "user", "post", "delete", false},

		// ── anonymous: no capabilities ──
		{"anonymous cannot create draft", "anonymous", "post:draft", "create", false},
		{"anonymous cannot publish", "anonymous", "post", "publish", false},
		{"anonymous cannot delete", "anonymous", "post", "delete", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := enforce(t, e, tt.sub, tt.obj, tt.act)
			if got != tt.expect {
				t.Errorf("Enforce(%q, %q, %q) = %v, want %v", tt.sub, tt.obj, tt.act, got, tt.expect)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════════════
// Test: Role inheritance chain
// ═══════════════════════════════════════════════════════════════════════

func TestRoleInheritance(t *testing.T) {
	e := setupTestEnforcer(t, defaultOpts())

	// super_admin inherits admin
	roles, err := e.GetImplicitRolesForUser("super_admin")
	if err != nil {
		t.Fatalf("GetImplicitRolesForUser error: %v", err)
	}
	assertContains(t, roles, "admin", "super_admin should inherit admin")
	assertContains(t, roles, "user", "super_admin should transitively inherit user")

	// admin inherits user
	roles, err = e.GetImplicitRolesForUser("admin")
	if err != nil {
		t.Fatalf("GetImplicitRolesForUser error: %v", err)
	}
	assertContains(t, roles, "user", "admin should inherit user")

	// user does NOT inherit admin
	roles, err = e.GetImplicitRolesForUser("user")
	if err != nil {
		t.Fatalf("GetImplicitRolesForUser error: %v", err)
	}
	assertNotContains(t, roles, "admin", "user must NOT inherit admin")
	assertNotContains(t, roles, "super_admin", "user must NOT inherit super_admin")
}

// ═══════════════════════════════════════════════════════════════════════
// Test: AdminCanDelete=false blocks route but capability still exists
// ═══════════════════════════════════════════════════════════════════════

func TestAdminCanDeleteFalse(t *testing.T) {
	opts := defaultOpts()
	opts.AdminCanDelete = false
	e := setupTestEnforcer(t, opts)

	// Route-level: admin cannot reach DELETE endpoint
	if enforce(t, e, "admin", "/api/v1/admin/posts/:id", "DELETE") {
		t.Error("admin should NOT be able to DELETE posts route when AdminCanDelete=false")
	}
	if enforce(t, e, "admin", "/api/v1/media/:id", "DELETE") {
		t.Error("admin should NOT be able to DELETE media route when AdminCanDelete=false")
	}
	if enforce(t, e, "admin", "/api/v1/tags/:id", "DELETE") {
		t.Error("admin should NOT be able to DELETE tags route when AdminCanDelete=false")
	}

	// Capability-level: the capability policy is always seeded (design note: mismatch)
	if !enforce(t, e, "admin", "post", "delete") {
		t.Error("admin should still have 'post delete' capability even when AdminCanDelete=false (route blocks it)")
	}

	// super_admin can still delete (bypass)
	if !enforce(t, e, "super_admin", "/api/v1/admin/posts/:id", "DELETE") {
		t.Error("super_admin should always be able to DELETE")
	}
}

// ═══════════════════════════════════════════════════════════════════════
// Test: AllowAnonymousRead=false removes anonymous policies
// ═══════════════════════════════════════════════════════════════════════

func TestAnonymousReadDisabled(t *testing.T) {
	opts := defaultOpts()
	opts.AllowAnonymousRead = false
	e := setupTestEnforcer(t, opts)

	// anonymous should be denied at policy level
	if enforce(t, e, "anonymous", "/api/v1/posts", "GET") {
		t.Error("anonymous should NOT be able to GET /api/v1/posts when AllowAnonymousRead=false")
	}
	if enforce(t, e, "anonymous", "/api/v1/posts/:id", "GET") {
		t.Error("anonymous should NOT be able to GET /api/v1/posts/:id when AllowAnonymousRead=false")
	}

	// NOTE: This test validates the POLICY layer. In production, the public post routes
	// (/api/v1/posts) are currently registered OUTSIDE the Casbin middleware group,
	// meaning this policy is never actually checked. See Issue #1 in the audit.
	// After the fix, these routes should go through Casbin and these policies will take effect.

	// authenticated users should still be able to read
	if !enforce(t, e, "user", "/api/v1/posts", "GET") {
		t.Error("user should still be able to GET /api/v1/posts")
	}
	if !enforce(t, e, "admin", "/api/v1/posts", "GET") {
		t.Error("admin should still be able to GET /api/v1/posts")
	}
}

// ═══════════════════════════════════════════════════════════════════════
// Test: UserCanUpload toggle
// ═══════════════════════════════════════════════════════════════════════

func TestUserUploadToggle(t *testing.T) {
	t.Run("upload disabled", func(t *testing.T) {
		opts := defaultOpts()
		opts.UserCanUpload = false
		e := setupTestEnforcer(t, opts)

		if enforce(t, e, "user", "/api/v1/media", "POST") {
			t.Error("user should NOT be able to POST media when UserCanUpload=false")
		}
		// read is always allowed
		if !enforce(t, e, "user", "/api/v1/media", "GET") {
			t.Error("user should be able to GET media")
		}
	})

	t.Run("upload enabled", func(t *testing.T) {
		opts := defaultOpts()
		opts.UserCanUpload = true
		e := setupTestEnforcer(t, opts)

		if !enforce(t, e, "user", "/api/v1/media", "POST") {
			t.Error("user should be able to POST media when UserCanUpload=true")
		}
	})
}

// ═══════════════════════════════════════════════════════════════════════
// Test: super_admin bypass does NOT depend on any explicit policy
// ═══════════════════════════════════════════════════════════════════════

func TestSuperAdminBypass_NoExplicitPolicy(t *testing.T) {
	opts := setupOpts{
		AdminCanDelete:     false,
		UserCanUpload:      false,
		AllowAnonymousRead: false,
	}
	e := setupTestEnforcer(t, opts)

	// Even without any super_admin policy, the matcher hardcodes the bypass
	tests := []struct {
		obj string
		act string
	}{
		{"/api/v1/admin/posts", "GET"},
		{"/api/v1/admin/posts/:id", "DELETE"},
		{"/api/v1/admin/posts/:id/publish", "POST"},
		{"/api/v1/anything", "PATCH"},
		{"post", "delete"},
		{"post", "publish"},
		{"post:draft", "create"},
	}

	for _, tt := range tests {
		t.Run(tt.obj+"_"+tt.act, func(t *testing.T) {
			if !enforce(t, e, "super_admin", tt.obj, tt.act) {
				t.Errorf("super_admin should bypass all checks for (%q, %q)", tt.obj, tt.act)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════════════
// Helpers
// ═══════════════════════════════════════════════════════════════════════

func assertContains(t *testing.T, slice []string, item string, msg string) {
	t.Helper()
	for _, s := range slice {
		if s == item {
			return
		}
	}
	t.Errorf("%s: %v does not contain %q", msg, slice, item)
}

func assertNotContains(t *testing.T, slice []string, item string, msg string) {
	t.Helper()
	for _, s := range slice {
		if s == item {
			t.Errorf("%s: %v should not contain %q", msg, slice, item)
			return
		}
	}
}
