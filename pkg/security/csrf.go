package security

import "github.com/google/uuid"

// GenerateToken 生成一个高熵随机字符串作为 CSRF 令牌
func GenerateToken() string {
	return uuid.New().String()
}
