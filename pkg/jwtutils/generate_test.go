package jwtutils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	t.Run("generates valid JWT", func(t *testing.T) {
		args := map[string]interface{}{
			"userId": "123",
			"role":   "moderator",
		}
		secret := "supersecret"
		expiration := time.Minute

		tokenStr, err := Generate(args, expiration, secret)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenStr)

		// Проверим, что токен можно распарсить и данные совпадают
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			assert.Equal(t, jwt.SigningMethodHS256, token.Method)
			return []byte(secret), nil
		})
		assert.NoError(t, err)
		assert.True(t, token.Valid)

		claims, ok := token.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, "123", claims["userId"])
		assert.Equal(t, "moderator", claims["role"])

		// Проверим, что exp установлен и валиден
		exp, ok := claims["exp"].(float64)
		assert.True(t, ok)
		assert.True(t, exp > float64(time.Now().Unix()))
	})

	t.Run("returns error if secret is empty", func(t *testing.T) {
		args := map[string]interface{}{"userId": "123"}
		tokenStr, err := Generate(args, time.Minute, "")
		assert.ErrorIs(t, err, jwt.ErrInvalidKey)
		assert.Empty(t, tokenStr)
	})
}
