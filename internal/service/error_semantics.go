package service

import (
	"KaldalisCMS/internal/core"
	"errors"
	"fmt"
)

func normalizeServiceError(err error) error {
	return normalizeServiceErrorWithOp("", err)
}

func normalizeServiceErrorWithOp(op string, err error) error {
	if err == nil {
		return nil
	}

	target := core.ErrInternalError
	switch {
	case errors.Is(err, core.ErrInvalidInput):
		target = core.ErrInvalidInput
	case errors.Is(err, core.ErrInvalidCredentials):
		target = core.ErrInvalidCredentials
	case errors.Is(err, core.ErrPermission):
		target = core.ErrPermission
	case errors.Is(err, core.ErrNotFound):
		target = core.ErrNotFound
	case errors.Is(err, core.ErrDuplicate):
		target = core.ErrDuplicate
	case errors.Is(err, core.ErrConflict):
		target = core.ErrConflict
	case errors.Is(err, core.ErrInternalError):
		target = core.ErrInternalError
	}

	if op != "" {
		return fmt.Errorf("%s: %w: %v", op, target, err)
	}
	return fmt.Errorf("%w: %v", target, err)
}

func normalizeServiceErrorWithOpMsg(op string, message string, err error) error {
	if err == nil {
		return nil
	}
	if message != "" {
		err = fmt.Errorf("%s: %w", message, err)
	}
	return normalizeServiceErrorWithOp(op, err)
}
