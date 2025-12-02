package core

import "errors"

// A list of common errors
var (
	ErrNotFound      = errors.New("not found")
	ErrDuplicate     = errors.New("duplicate entry")
	ErrInvalidInput  = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDBConnection  = errors.New("database connection error")
	ErrTransaction   = errors.New("database transaction error")
	ErrPermission    = errors.New("permission denied")
	ErrInternalError = errors.New("internal server error") // General purpose internal error
)
