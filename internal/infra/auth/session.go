package auth

import (
	"net/http"
	"time"

	"KaldalisCMS/pkg/auth"
	"KaldalisCMS/pkg/security"

	"github.com/spf13/viper"
)

type Config struct {
	Secret     []byte // 业务直接用 bytes
	TTL        time.Duration
	AuthCookie string
	CSRFCookie string
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
		Path:       raw.Path,
		Domain:     raw.Domain,
		Secure:     raw.Secure,
		SameSite:   http.SameSite(raw.SameSite),
	}, nil
}

// EstablishSession 封装了登录时同时设置 JWT 和 CSRF Cookie 的逻辑
func EstablishSession(w http.ResponseWriter, cfg Config, userID uint) error {
	token, err := auth.Generate(userID, cfg.Secret, cfg.TTL)
	if err != nil {
		return err
	}
	csrf := security.GenerateToken()

	// Auth Cookie (HttpOnly)
	setCookie(w, cfg, cfg.AuthCookie, token, true)
	// CSRF Cookie
	setCookie(w, cfg, cfg.CSRFCookie, csrf, false)
	return nil
}

func DestroySession(w http.ResponseWriter, cfg Config) {
	deleteCookie(w, cfg, cfg.AuthCookie)
	deleteCookie(w, cfg, cfg.CSRFCookie)
}

func setCookie(w http.ResponseWriter, cfg Config, name, value string, httpOnly bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     cfg.Path,
		Domain:   cfg.Domain,
		MaxAge:   int(cfg.TTL.Seconds()),
		HttpOnly: httpOnly,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
	})
}

func deleteCookie(w http.ResponseWriter, cfg Config, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     cfg.Path,
		Domain:   cfg.Domain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
	})
}
