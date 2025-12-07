package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT密钥（学习阶段可写死，生产环境请用环境变量）
var jwtSecret = []byte("my_secret_key")

// MakeToken 生成 JWT token，并返回 token 字符串
func MakeToken(username string, expire time.Time) (string, time.Time, error) {
	// 创建 payload
	claims := jwt.MapClaims{
		"username": username,
		"exp":      expire.Unix(),
		"iat":      time.Now().Unix(),
	}

	// 使用 HMAC SHA256 签名生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expire, nil
}

// ParseToken 解析 JWT token，并返回用户名和过期时间
func ParseToken(tokenStr string) (string, time.Time, error) {
	if tokenStr == "" {
		return "", time.Time{}, errors.New("token is empty")
	}

	claims := jwt.MapClaims{}

	// 解析 token
	t, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", time.Time{}, err
	}

	if !t.Valid {
		return "", time.Time{}, errors.New("invalid token")
	}

	// 获取 username
	username, ok := claims["username"].(string)
	if !ok {
		return "", time.Time{}, errors.New("username not found in token")
	}

	// 获取 exp
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return "", time.Time{}, errors.New("exp not found in token")
	}
	expireTime := time.Unix(int64(expFloat), 0)

	return username, expireTime, nil
}

// MakeRefreshToken 生成 refresh token（简单版）
func MakeRefreshToken(username string) (string, time.Time, error) {
	// refresh token 有效期比 access token 长
	expire := time.Now().Add(24 * time.Hour)
	return MakeToken(username, expire)
}
