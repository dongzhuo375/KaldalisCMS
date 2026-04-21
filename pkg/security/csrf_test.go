package security

import (
	"regexp"
	"testing"
)

// uuidV4 format: 8-4-4-4-12 hex chars, with version 4 and variant bits.
var uuidV4Re = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestGenerateToken(t *testing.T) {
	t.Run("is non-empty and UUIDv4 shaped", func(t *testing.T) {
		tok := GenerateToken()
		if !uuidV4Re.MatchString(tok) {
			t.Errorf("token %q is not a UUIDv4", tok)
		}
	})

	t.Run("unique across many calls", func(t *testing.T) {
		const n = 1000
		seen := make(map[string]struct{}, n)
		for i := range n {
			tok := GenerateToken()
			if _, dup := seen[tok]; dup {
				t.Fatalf("collision after %d tokens: %q", i, tok)
			}
			seen[tok] = struct{}{}
		}
	})
}
