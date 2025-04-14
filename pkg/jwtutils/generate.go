package jwtutils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func Generate(args map[string]interface{}, expiration time.Duration, secret string) (string, error) {
	if secret == "" {
		return "", jwt.ErrInvalidKey
	}

	claims := jwt.MapClaims(args)
	claims["exp"] = time.Now().Add(expiration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
