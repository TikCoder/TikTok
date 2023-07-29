package jwtUtils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"user_name"`
	jwt.RegisteredClaims
}

var Secret = []byte("tiktok")

// GenToken 生成 token
func GenToken(userId int64, userName string) (string, error) {
	claims := JWTClaims{
		Username: userName,
		UserId:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "service",
		},
	}
	// SigningMethodHS256 不是 SigningMethodES256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(Secret)
	if err != nil {
		// todo log err
		return "", err
	}
	return signedToken, nil
}

// ParseToken 解析 token
func ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// VerifyToken 验证 token
func VerifyToken(tokenString string) (int64, error) {

	if tokenString == "" {
		return int64(0), nil
	}
	claims, err := ParseToken(tokenString)
	if err != nil {
		// todo : log err
		return int64(0), err
	}
	return claims.UserId, nil
}
