package errorx

import (
	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/core"

	"github.com/gin-gonic/gin"
)

const (
	CtxRequestIDKey = "kaldalis_request_id"
	CtxErrorCodeKey = "kaldalis_error_code"
	HeaderRequestID = "X-Request-Id"
)

func RespondMessage(c *gin.Context, status int, message string) {
	c.JSON(status, dto.MessageResponse{Message: message})
}

func RespondError(c *gin.Context, status int, code core.ErrorCode, message string, details map[string]any) {
	writeError(c, status, code, message, details)
}

func AbortError(c *gin.Context, status int, code core.ErrorCode, message string, details map[string]any) {
	writeError(c, status, code, message, details)
	c.Abort()
}

func writeError(c *gin.Context, status int, code core.ErrorCode, message string, details map[string]any) {
	if status <= 0 {
		status = core.HTTPStatusOf(code)
	}
	if message == "" {
		message = core.DefaultMessageOf(code)
	}
	details = withRequestID(c, details)
	details = core.SanitizeDetails(code, details)

	c.Set(CtxErrorCodeKey, string(code))
	c.JSON(status, dto.ErrorResponse{
		Code:    string(code),
		Message: message,
		Details: details,
	})
}

func RespondErrorByCore(c *gin.Context, err error, defaultStatus int, details map[string]any) {
	code := core.ErrorCodeOf(err)
	_ = defaultStatus // kept for backward compatibility at call sites
	RespondError(c, core.HTTPStatusOf(code), code, core.DefaultMessageOf(code), details)
}

func RespondValidationError(c *gin.Context, message string, details map[string]any) {
	if message == "" {
		message = core.DefaultMessageOf(core.CodeValidationFailed)
	}
	RespondError(c, core.HTTPStatusOf(core.CodeValidationFailed), core.CodeValidationFailed, message, details)
}

func RespondTimeoutError(c *gin.Context, message string) {
	if message == "" {
		message = core.DefaultMessageOf(core.CodeTimeout)
	}
	RespondError(c, core.HTTPStatusOf(core.CodeTimeout), core.CodeTimeout, message, nil)
}

func RespondInternalError(c *gin.Context) {
	RespondError(c, core.HTTPStatusOf(core.CodeInternalError), core.CodeInternalError, core.DefaultMessageOf(core.CodeInternalError), nil)
}

func AbortInternalError(c *gin.Context) {
	AbortError(c, core.HTTPStatusOf(core.CodeInternalError), core.CodeInternalError, core.DefaultMessageOf(core.CodeInternalError), nil)
}

func withRequestID(c *gin.Context, details map[string]any) map[string]any {
	if details == nil {
		details = map[string]any{}
	}
	if rid, ok := c.Get(CtxRequestIDKey); ok {
		if requestID, ok := rid.(string); ok && requestID != "" {
			details["request_id"] = requestID
			c.Writer.Header().Set(HeaderRequestID, requestID)
		}
	}
	return details
}
