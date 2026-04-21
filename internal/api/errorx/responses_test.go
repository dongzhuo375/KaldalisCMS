package errorx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"KaldalisCMS/internal/api/v1/dto"
	"KaldalisCMS/internal/core"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestCtx() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, w
}

func decodeErr(t *testing.T, body []byte) dto.ErrorResponse {
	t.Helper()
	var got dto.ErrorResponse
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("decode: %v (body=%s)", err, body)
	}
	return got
}

func TestRespondErrorByCore_Mapping(t *testing.T) {
	// SanitizeDetails enforces a per-code allow-list: `field` is only kept for
	// validation + duplicate. Every case passes the same payload and asserts
	// survival based on that policy, so this test pins both the code→HTTP
	// mapping and the sanitization contract.
	cases := []struct {
		name      string
		err       error
		wantCode  core.ErrorCode
		wantHTTP  int
		fieldKept bool
	}{
		{"validation", core.ErrInvalidInput, core.CodeValidationFailed, http.StatusBadRequest, true},
		{"credentials", core.ErrInvalidCredentials, core.CodeUnauthorized, http.StatusUnauthorized, false},
		{"forbidden", core.ErrPermission, core.CodeForbidden, http.StatusForbidden, false},
		{"not found", core.ErrNotFound, core.CodeNotFound, http.StatusNotFound, false},
		{"duplicate", core.ErrDuplicate, core.CodeDuplicateResource, http.StatusConflict, true},
		{"conflict", core.ErrConflict, core.CodeConflict, http.StatusConflict, false},
		{"unknown -> internal", errors.New("unexpected"), core.CodeInternalError, http.StatusInternalServerError, false},
		{"wrapped not found", fmt.Errorf("lookup: %w", core.ErrNotFound), core.CodeNotFound, http.StatusNotFound, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, w := newTestCtx()
			RespondErrorByCore(c, tc.err, http.StatusTeapot, map[string]any{"field": "title"})
			if w.Code != tc.wantHTTP {
				t.Fatalf("http: want %d got %d", tc.wantHTTP, w.Code)
			}
			got := decodeErr(t, w.Body.Bytes())
			if got.Code != string(tc.wantCode) {
				t.Fatalf("code: want %s got %s", tc.wantCode, got.Code)
			}
			if got.Message == "" {
				t.Fatal("message empty")
			}
			if tc.fieldKept {
				if got.Details["field"] != "title" {
					t.Fatalf("field should survive for %s: %+v", tc.wantCode, got.Details)
				}
			} else if _, present := got.Details["field"]; present {
				t.Fatalf("field should be sanitized for %s: %+v", tc.wantCode, got.Details)
			}
		})
	}
}

func TestSanitizeDetails_StripsSensitiveKeys(t *testing.T) {
	c, w := newTestCtx()
	RespondError(c, http.StatusBadRequest, core.CodeValidationFailed, "bad", map[string]any{
		"field":         "title",
		"password":      "hunter2",
		"auth_token":    "abc",
		"Authorization": "Bearer x",
	})
	got := decodeErr(t, w.Body.Bytes())
	if got.Details["field"] != "title" {
		t.Fatalf("allowed key dropped: %+v", got.Details)
	}
	for _, bad := range []string{"password", "auth_token", "Authorization"} {
		if _, present := got.Details[bad]; present {
			t.Fatalf("sensitive key %q leaked: %+v", bad, got.Details)
		}
	}
}

func TestRespondValidationError_DefaultMessage(t *testing.T) {
	c, w := newTestCtx()
	RespondValidationError(c, "", nil)
	got := decodeErr(t, w.Body.Bytes())
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", w.Code)
	}
	if got.Code != string(core.CodeValidationFailed) || got.Message == "" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestRespondTimeoutError(t *testing.T) {
	c, w := newTestCtx()
	RespondTimeoutError(c, "slow upstream")
	if w.Code != http.StatusGatewayTimeout {
		t.Fatalf("status: %d", w.Code)
	}
	got := decodeErr(t, w.Body.Bytes())
	if got.Code != string(core.CodeTimeout) || got.Message != "slow upstream" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestRespondInternalError(t *testing.T) {
	c, w := newTestCtx()
	RespondInternalError(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", w.Code)
	}
	got := decodeErr(t, w.Body.Bytes())
	if got.Code != string(core.CodeInternalError) {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestRespondMessage(t *testing.T) {
	c, w := newTestCtx()
	RespondMessage(c, http.StatusCreated, "ok")
	if w.Code != http.StatusCreated {
		t.Fatalf("status: %d", w.Code)
	}
	var got dto.MessageResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Message != "ok" {
		t.Fatalf("message: %q", got.Message)
	}
}
