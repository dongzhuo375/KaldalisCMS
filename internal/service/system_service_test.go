package service

import (
	"context"
	"errors"
	"testing"

	"KaldalisCMS/internal/core"
)

// SetupOnce validates its params before touching the database, so these branches
// can be covered without any DB wiring. Full-path tests (lock, mark-installed, create-admin)
// belong in an integration suite that stands up a real Postgres or sqlite fixture.

func TestSystemService_SetupOnce_ValidationErrors(t *testing.T) {
	ctx := context.Background()
	svc := &SystemService{} // DB deps unused on this path

	cases := []struct {
		name string
		p    SetupParams
	}{
		{"empty site", SetupParams{AdminUsername: "a", AdminEmail: "a@b", AdminPassword: "p"}},
		{"empty username", SetupParams{SiteName: "s", AdminEmail: "a@b", AdminPassword: "p"}},
		{"empty email", SetupParams{SiteName: "s", AdminUsername: "a", AdminPassword: "p"}},
		{"empty password", SetupParams{SiteName: "s", AdminUsername: "a", AdminEmail: "a@b"}},
		{"whitespace stripped to empty", SetupParams{SiteName: "   ", AdminUsername: "a", AdminEmail: "a@b", AdminPassword: "p"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := svc.SetupOnce(ctx, tc.p)
			if !errors.Is(err, core.ErrInvalidInput) {
				t.Fatalf("want ErrInvalidInput, got %v", err)
			}
		})
	}
}
