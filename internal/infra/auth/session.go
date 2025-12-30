package auth

import (
	"errors"
	"net/http"
	"time"

	"KaldalisCMS/pkg/auth"
	"KaldalisCMS/pkg/security"

	"github.com/spf13/viper"
)

var ErrNoToken = errors.New("no token found")

type Config struct {
	Secret     []byte // 业务直接用 bytes
	TTL        time.Duration
	AuthCookie string
	CSRFCookie string
	RoleCookie string
	Path       string
	Domain     string
	Secure     bool
	SameSite   http.SameSite // 业务直接用强类型
}

// 内部使用的解析结构（私有）
type rawConfig struct {
	Secret     string        `mapstructure:"secret"`
	TTL        time.Duration `mapstructure:"ttl"`
	AuthCookie string        `mapstructure:"auth_cookie"`
	CSRFCookie string        `mapstructure:"csrf_cookie"`
	RoleCookie string 		 `mapstructure:"role_cookie"`
	Path       string        `mapstructure:"path"`
	Domain     string        `mapstructure:"domain"`
	Secure     bool          `mapstructure:"secure"`
	SameSite   int           `mapstructure:"same_site"`
}

// LoadConfig 负责转换业务config
func LoadConfig(v *viper.Viper) (*Config, error) {
	var raw rawConfig

	sub := v.Sub("jwt")

	if err := sub.Unmarshal(&raw); err != nil {
		return nil, err
	}

	return &Config{
		Secret:     []byte(raw.Secret),
		TTL:        raw.TTL,
		AuthCookie: raw.AuthCookie,
		CSRFCookie: raw.CSRFCookie,
		RoleCookie: raw.RoleCookie,
		Path:       raw.Path,
		Domain:     raw.Domain,
		Secure:     raw.Secure,
		SameSite:   http.SameSite(raw.SameSite),
	}, nil
}

type SessionManager struct {
	cfg Config // 配置已经被锁死在实例里了
}

func NewSessionManager(cfg Config) *SessionManager {
	return &SessionManager{cfg: cfg}
}

// EstablishSession 封装了登录时同时设置 JWT 和 CSRF Cookie 的逻辑
func (m *SessionManager) EstablishSession(w http.ResponseWriter, userID uint ,role string) error {
	csrf := security.GenerateToken()
	token, err := auth.GenerateHashCSRF(userID, role, m.cfg.Secret, m.cfg.TTL, csrf)

	if err != nil {
		return err
	}
	// Auth Cookie (HttpOnly)
	m.setCookie(w, m.cfg.AuthCookie, token, true)
	//Role cookie
	m.setCookie(w, m.cfg.RoleCookie, role, false)
	// CSRF Cookie
	m.setCookie(w, m.cfg.CSRFCookie, csrf, false)
	return nil
}

func (m *SessionManager) DestroySession(w http.ResponseWriter) {
	m.deleteCookie(w, m.cfg.AuthCookie)
	m.deleteCookie(w, m.cfg.CSRFCookie)
	m.deleteCookie(w, m.cfg.RoleCookie)
}

func (m *SessionManager) setCookie(w http.ResponseWriter, name, value string, httpOnly bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     m.cfg.Path,
		Domain:   m.cfg.Domain,
		MaxAge:   int(m.cfg.TTL.Seconds()),
		HttpOnly: httpOnly,
		Secure:   m.cfg.Secure,
		SameSite: m.cfg.SameSite,
	})
}

func (m *SessionManager) deleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     m.cfg.Path,
		Domain:   m.cfg.Domain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   m.cfg.Secure,
		SameSite: m.cfg.SameSite,
	})
}

// 从请求中提取并校验身份
func (m *SessionManager) Authenticate(r *http.Request) (*auth.CustomClaims, error) {
	token := ""

	//尝试从 Cookie 获取
	if ck, err := r.Cookie(m.cfg.AuthCookie); err == nil {
		token = ck.Value
	}

	//尝试从 Header 获取
	if token == "" {
		ah := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if len(ah) > len(prefix) && ah[:len(prefix)] == prefix {
			token = ah[len(prefix):]
		}
	}

	if token == "" {
		return nil, ErrNoToken
	}

	// 校验 Token（使用内部持有的 Secret）
	return auth.Parse(token, m.cfg.Secret)
}

func (m *SessionManager) ValidateCSRF(r *http.Request, expectedHash string) error {
	//Cookie
	cookie, err := r.Cookie(m.cfg.CSRFCookie)
	if err != nil {
		return errors.New("CSRF cookie missing")
	}

	//Header
	headerVal := r.Header.Get("X-CSRF-Token")
	if headerVal == "" || headerVal != cookie.Value {
		return errors.New("CSRF token mismatch")
	}

	//校验指纹绑定
	if expectedHash != "" {
		if auth.HashToken(headerVal) != expectedHash {
			return errors.New("CSRF token binding invalid")
		}
	}

	return nil
}

func (m *SessionManager) GetTTL() time.Duration {
	return m.cfg.TTL
}
