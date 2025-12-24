package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID uint   `json:"user_id"`
	CsrfH  string `json:"csrf_h"` // 新增：用于存储 CSRF 指纹
	jwt.RegisteredClaims
}

// Generate 生成 JWT 字符串，保留
//func Generate(userID uint, secret []byte, ttl time.Duration) (string, error) {
//	claims := CustomClaims{
//		UserID: userID,
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
//			IssuedAt:  jwt.NewNumericDate(time.Now()),
//			Issuer:    "KaldalisCMS",
//		},
//	}
//	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
//}

// Parse 解析并验证 JWT 字符串
func Parse(tokenStr string, secret []byte) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid // 使用标准库提供的错误码
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid claims")
}

// 生成绑定CSRF的JWT字符串
func GenerateHashCSRF(userID uint, secret []byte, ttl time.Duration, csrfToken string) (string, error) {
	csrfHash := HashToken(csrfToken)

	claims := CustomClaims{
		UserID: userID,
		CsrfH:  csrfHash, // 绑定指纹
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "KaldalisCMS",
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func HashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}
