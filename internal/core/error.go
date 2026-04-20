package core

import (
	"errors"
	"net/http"
	"strings"
)

// ErrorCode is a stable machine-readable API error code.
type ErrorCode string

const (
	CodeValidationFailed  ErrorCode = "VALIDATION_FAILED"
	CodeUnauthorized      ErrorCode = "UNAUTHORIZED"
	CodeForbidden         ErrorCode = "FORBIDDEN"
	CodeNotFound          ErrorCode = "NOT_FOUND"
	CodeDuplicateResource ErrorCode = "DUPLICATE_RESOURCE"
	CodeConflict          ErrorCode = "CONFLICT"
	CodeTimeout           ErrorCode = "TIMEOUT"
	CodeInternalError     ErrorCode = "INTERNAL_ERROR"
)

// ErrorPolicy defines how an error code is exposed over HTTP.
type ErrorPolicy struct {
	HTTPStatus      int
	Message         string
	AllowDetailsKey map[string]struct{}
}

var errorPolicyCatalog = map[ErrorCode]ErrorPolicy{
	CodeValidationFailed: {
		HTTPStatus: http.StatusBadRequest,
		Message:    "request validation failed",
		AllowDetailsKey: map[string]struct{}{
			"field":      {},
			"fields":     {},
			"reason":     {},
			"request_id": {},
		},
	},
	CodeUnauthorized: {
		HTTPStatus: http.StatusUnauthorized,
		Message:    "unauthorized",
		AllowDetailsKey: map[string]struct{}{
			"request_id": {},
		},
	},
	CodeForbidden: {
		HTTPStatus: http.StatusForbidden,
		Message:    "permission denied",
		AllowDetailsKey: map[string]struct{}{
			"request_id": {},
		},
	},
	CodeNotFound: {
		HTTPStatus: http.StatusNotFound,
		Message:    "resource not found",
		AllowDetailsKey: map[string]struct{}{
			"resource":   {},
			"id":         {},
			"request_id": {},
		},
	},
	CodeDuplicateResource: {
		HTTPStatus: http.StatusConflict,
		Message:    "resource already exists",
		AllowDetailsKey: map[string]struct{}{
			"field":      {},
			"request_id": {},
		},
	},
	CodeConflict: {
		HTTPStatus: http.StatusConflict,
		Message:    "request conflict",
		AllowDetailsKey: map[string]struct{}{
			"resource":   {},
			"references": {},
			"request_id": {},
		},
	},
	CodeTimeout: {
		HTTPStatus: http.StatusGatewayTimeout,
		Message:    "request timed out",
		AllowDetailsKey: map[string]struct{}{
			"request_id": {},
		},
	},
	CodeInternalError: {
		HTTPStatus: http.StatusInternalServerError,
		Message:    "internal server error",
		AllowDetailsKey: map[string]struct{}{
			"request_id": {},
		},
	},
}

// A list of common errors
var (
	ErrNotFound           = errors.New("not found")
	ErrDuplicate          = errors.New("duplicate entry")
	ErrConflict           = errors.New("conflict")
	ErrTimeout            = errors.New("timeout")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDBConnection       = errors.New("database connection error")
	ErrTransaction        = errors.New("database transaction error")
	ErrPermission         = errors.New("permission denied")
	ErrInternalError      = errors.New("internal server error") // General purpose internal error
)

// ErrorCodeOf maps domain errors to stable API codes.
func ErrorCodeOf(err error) ErrorCode {
	switch {
	case errors.Is(err, ErrInvalidInput):
		return CodeValidationFailed
	case errors.Is(err, ErrInvalidCredentials):
		return CodeUnauthorized
	case errors.Is(err, ErrPermission):
		return CodeForbidden
	case errors.Is(err, ErrNotFound):
		return CodeNotFound
	case errors.Is(err, ErrDuplicate):
		return CodeDuplicateResource
	case errors.Is(err, ErrConflict):
		return CodeConflict
	case errors.Is(err, ErrTimeout):
		return CodeTimeout
	default:
		return CodeInternalError
	}
}

func ErrorPolicyOf(code ErrorCode) ErrorPolicy {
	if p, ok := errorPolicyCatalog[code]; ok {
		return p
	}
	return errorPolicyCatalog[CodeInternalError]
}

func HTTPStatusOf(code ErrorCode) int {
	return ErrorPolicyOf(code).HTTPStatus
}

func DefaultMessageOf(code ErrorCode) string {
	return ErrorPolicyOf(code).Message
}

// SanitizeDetails enforces a strict allow-list by error code and strips likely secret keys.
func SanitizeDetails(code ErrorCode, details map[string]any) map[string]any {
	if details == nil {
		return map[string]any{}
	}

	policy := ErrorPolicyOf(code)
	out := make(map[string]any)
	for key, value := range details {
		if isSensitiveDetailKey(key) {
			continue
		}
		if _, ok := policy.AllowDetailsKey[key]; !ok {
			continue
		}
		out[key] = value
	}
	if out == nil {
		return map[string]any{}
	}
	return out
}

func isSensitiveDetailKey(key string) bool {
	lower := strings.ToLower(key)
	return strings.Contains(lower, "password") ||
		strings.Contains(lower, "token") ||
		strings.Contains(lower, "secret") ||
		strings.Contains(lower, "authorization")
}
