package core

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorCodeOf(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want ErrorCode
	}{
		{"invalid input", ErrInvalidInput, CodeValidationFailed},
		{"invalid credentials", ErrInvalidCredentials, CodeUnauthorized},
		{"permission denied", ErrPermission, CodeForbidden},
		{"not found", ErrNotFound, CodeNotFound},
		{"duplicate", ErrDuplicate, CodeDuplicateResource},
		{"conflict", ErrConflict, CodeConflict},

		{"wrapped not found", fmt.Errorf("load post: %w", ErrNotFound), CodeNotFound},
		{"wrapped duplicate", fmt.Errorf("save user: %w", ErrDuplicate), CodeDuplicateResource},

		{"unknown error falls back to internal", errors.New("some random failure"), CodeInternalError},
		{"db connection falls back to internal", ErrDBConnection, CodeInternalError},
		{"transaction falls back to internal", ErrTransaction, CodeInternalError},
		{"generic internal", ErrInternalError, CodeInternalError},

		{"nil error falls back to internal", nil, CodeInternalError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ErrorCodeOf(tt.err)
			if got != tt.want {
				t.Errorf("ErrorCodeOf(%v) = %q, want %q", tt.err, got, tt.want)
			}
		})
	}
}
