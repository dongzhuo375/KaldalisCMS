package v1

import (
	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

func respondMessage(c *gin.Context, status int, message string) {
	c.JSON(status, dto.MessageResponse{Message: message})
}

func respondError(c *gin.Context, status int, code core.ErrorCode, message string, details map[string]any) {
	c.JSON(status, dto.ErrorResponse{
		Code:    string(code),
		Message: message,
		Details: details,
	})
}

func respondErrorByCore(c *gin.Context, err error, defaultStatus int, details map[string]any) {
	code := core.ErrorCodeOf(err)
	status := defaultStatus
	message := "internal server error"

	switch code {
	case core.CodeValidationFailed:
		status = http.StatusBadRequest
		message = "request validation failed"
	case core.CodeUnauthorized:
		status = http.StatusUnauthorized
		message = "unauthorized"
	case core.CodeForbidden:
		status = http.StatusForbidden
		message = "permission denied"
	case core.CodeNotFound:
		status = http.StatusNotFound
		message = "resource not found"
	case core.CodeDuplicateResource:
		status = http.StatusConflict
		message = "resource already exists"
	case core.CodeConflict:
		status = http.StatusConflict
		message = "request conflict"
	case core.CodeTimeout:
		status = http.StatusGatewayTimeout
		message = "request timed out"
	default:
		status = http.StatusInternalServerError
	}

	respondError(c, status, code, message, details)
}

func respondValidationError(c *gin.Context, message string, details map[string]any) {
	if message == "" {
		message = "request validation failed"
	}
	respondError(c, http.StatusBadRequest, core.CodeValidationFailed, message, details)
}

func respondTimeoutError(c *gin.Context, message string) {
	if message == "" {
		message = "request timed out"
	}
	respondError(c, http.StatusGatewayTimeout, core.CodeTimeout, message, nil)
}

func respondInternalError(c *gin.Context) {
	respondError(c, http.StatusInternalServerError, core.CodeInternalError, "internal server error", nil)
}
