package core

import "errors"

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

// A list of common errors
var (
	ErrNotFound           = errors.New("not found")
	ErrDuplicate          = errors.New("duplicate entry")
	ErrConflict           = errors.New("conflict")
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
	default:
		return CodeInternalError
	}
}
