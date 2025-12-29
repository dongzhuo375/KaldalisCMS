package core

import (
	"KaldalisCMS/pkg/auth"
	"net/http"
	"time"
)

// SeesionManager接口定义
type SessionManager interface {
	EstablishSession(w http.ResponseWriter, userID uint, role string) error
	DestroySession(w http.ResponseWriter)
	Authenticate(r *http.Request) (*auth.CustomClaims, error)
	ValidateCSRF(r *http.Request, expectedHash string) error
	GetTTL() time.Duration
}
