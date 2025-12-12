package core

import "net/http"

// AuthManager 是 core 层对认证 infra 的抽象（用于依赖倒置）
type AuthManager interface {
	// Login 签发 token 并写入 auth + csrf cookie
	Login(w http.ResponseWriter, userID uint, secureFlag bool) error

	// Logout 清除 cookie
	Logout(w http.ResponseWriter)

	// Parse 解析 token（从中间件/其他地方需要直接解析时使用）
	Parse(tokenStr string) (map[string]interface{}, error)

	// 返回配置（可选），例如 auth cookie 名称 / csrf cookie 名
	AuthCookieName() string
	CSRFCookieName() string
}
