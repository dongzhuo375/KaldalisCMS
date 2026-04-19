package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var testSecret = []byte("test-secret-key-do-not-use-in-prod")

func TestHashToken(t *testing.T) {
	t.Run("deterministic", func(t *testing.T) {
		h1 := HashToken("abc")
		h2 := HashToken("abc")
		if h1 != h2 {
			t.Errorf("HashToken not deterministic: %q vs %q", h1, h2)
		}
	})

	t.Run("different input yields different hash", func(t *testing.T) {
		if HashToken("abc") == HashToken("abd") {
			t.Error("different inputs produced the same hash")
		}
	})

	t.Run("empty string is hashable", func(t *testing.T) {
		h := HashToken("")
		// sha256("") is a well-known constant
		const emptySha256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		if h != emptySha256 {
			t.Errorf("HashToken(\"\") = %q, want %q", h, emptySha256)
		}
	})

	t.Run("output is hex-encoded 64 chars", func(t *testing.T) {
		h := HashToken("anything")
		if len(h) != 64 {
			t.Errorf("expected 64-char hex, got %d chars", len(h))
		}
	})
}

func TestGenerateAndParse_RoundTrip(t *testing.T) {
	const userID uint = 42
	const role = "admin"
	const csrf = "csrf-plaintext-123"

	tokenStr, err := GenerateHashCSRF(userID, role, testSecret, time.Hour, csrf)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := Parse(tokenStr, testSecret)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID = %d, want %d", claims.UserID, userID)
	}
	if claims.Role != role {
		t.Errorf("Role = %q, want %q", claims.Role, role)
	}
	if claims.CsrfH != HashToken(csrf) {
		t.Errorf("CsrfH not bound to the csrf token fingerprint")
	}
	if claims.Issuer != "KaldalisCMS" {
		t.Errorf("Issuer = %q, want KaldalisCMS", claims.Issuer)
	}
}

func TestParse_WrongSecret(t *testing.T) {
	tokenStr, err := GenerateHashCSRF(1, "user", testSecret, time.Hour, "csrf")
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	_, err = Parse(tokenStr, []byte("different-secret"))
	if err == nil {
		t.Fatal("expected error when parsing with wrong secret, got nil")
	}
}

func TestParse_TamperedToken(t *testing.T) {
	tokenStr, err := GenerateHashCSRF(1, "user", testSecret, time.Hour, "csrf")
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	// Flip a character in the payload section (middle segment).
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3 JWT parts, got %d", len(parts))
	}
	tampered := parts[0] + "." + parts[1] + "X." + parts[2]

	if _, err := Parse(tampered, testSecret); err == nil {
		t.Fatal("expected error for tampered token, got nil")
	}
}

func TestParse_ExpiredToken(t *testing.T) {
	tokenStr, err := GenerateHashCSRF(1, "user", testSecret, -time.Minute, "csrf")
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	_, err = Parse(tokenStr, testSecret)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestParse_RejectsNoneAlgorithm(t *testing.T) {
	// Craft a token using alg=none — a classic JWT attack vector.
	claims := CustomClaims{
		UserID: 1,
		Role:   "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenStr, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("signing with none failed: %v", err)
	}

	if _, err := Parse(tokenStr, testSecret); err == nil {
		t.Fatal("expected error for alg=none token, got nil")
	}
}

func TestParse_MalformedToken(t *testing.T) {
	cases := []string{"", "not-a-jwt", "a.b", "a.b.c.d"}
	for _, s := range cases {
		t.Run(s, func(t *testing.T) {
			if _, err := Parse(s, testSecret); err == nil {
				t.Errorf("expected error for %q, got nil", s)
			}
		})
	}
}
