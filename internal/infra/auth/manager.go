package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

//预留，cookie名、secure、sameSite从硬编码改为配置中读取

// Manager 管理 JWT/Cookie/CSRF 行为
type Manager struct {
	Secret        string
	TTL           time.Duration
	AuthCookie    string
	CSRFCookie    string
	Path          string
	Domain        string
	Secure        bool
	SameSite      http.SameSite
	CSRFCookieLen int // bytes, default 32
}

// NewManager 创建 Manager
func NewManager(secret string, ttl time.Duration, authCookie, csrfCookie, path, domain string, secure bool, sameSite http.SameSite) *Manager {
	return &Manager{
		Secret:        secret,
		TTL:           ttl,
		AuthCookie:    authCookie,
		CSRFCookie:    csrfCookie,
		Path:          path,
		Domain:        domain,
		Secure:        secure,
		SameSite:      sameSite,
		CSRFCookieLen: 32,
	}
}

// Login 为 userID 签发 JWT，并写入 auth + csrf cookie（secureFlag 可由外层基于 TLS 或配置决定）
func (m *Manager) Login(w http.ResponseWriter, userID uint, secureFlag bool) error {
	token, err := GenerateToken(userID, m.Secret, m.TTL)
	if err != nil {
		return err
	}

	expiry := time.Now().Add(m.TTL)

	// auth cookie (HttpOnly)
	http.SetCookie(w, &http.Cookie{
		Name:     m.AuthCookie,
		Value:    token,
		Path:     m.Path,
		Domain:   m.Domain,
		HttpOnly: true,
		Secure:   secureFlag && m.Secure,
		SameSite: m.SameSite,
		Expires:  expiry,
	})

	// csrf token (readable by JS) - double submit
	csrfBytes := make([]byte, m.CSRFCookieLen)
	if _, err := rand.Read(csrfBytes); err != nil {
		// 清理刚刚设置的 auth cookie，防止残留
		http.SetCookie(w, &http.Cookie{
			Name:     m.AuthCookie,
			Value:    "",
			Path:     m.Path,
			Domain:   m.Domain,
			HttpOnly: true,
			Secure:   secureFlag && m.Secure,
			SameSite: m.SameSite,
			MaxAge:   -1,
		})
		return err
	}
	csrfToken := base64.RawURLEncoding.EncodeToString(csrfBytes)
	http.SetCookie(w, &http.Cookie{
		Name:     m.CSRFCookie,
		Value:    csrfToken,
		Path:     m.Path,
		Domain:   m.Domain,
		HttpOnly: false, // front-end reads this
		Secure:   secureFlag && m.Secure,
		SameSite: m.SameSite,
		Expires:  expiry,
	})

	return nil
}

// Logout 清除 cookie
func (m *Manager) Logout(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     m.AuthCookie,
		Value:    "",
		Path:     m.Path,
		Domain:   m.Domain,
		HttpOnly: true,
		Secure:   m.Secure,
		SameSite: m.SameSite,
		MaxAge:   -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     m.CSRFCookie,
		Value:    "",
		Path:     m.Path,
		Domain:   m.Domain,
		HttpOnly: false,
		Secure:   m.Secure,
		SameSite: m.SameSite,
		MaxAge:   -1,
	})
}

// Parse parses token string and returns claims map
func (m *Manager) Parse(tokenStr string) (map[string]interface{}, error) {
	claims, err := ParseToken(tokenStr, m.Secret)
	if err != nil {
		return nil, err
	}
	// convert jwt.MapClaims to map[string]interface{}
	res := make(map[string]interface{})
	for k, v := range claims {
		res[k] = v
	}
	return res, nil
}

func (m *Manager) AuthCookieName() string {
	return m.AuthCookie
}
func (m *Manager) CSRFCookieName() string {
	return m.CSRFCookie
}
