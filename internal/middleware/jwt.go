package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"pvz-service/internal/handler/pkg/response"
)

const ErrInvalidToken = "Invalid token"

const (
	UserIDKey = "userId"
	RoleKey   = "role"
)

type JWT struct {
	secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{secret}
}
func (j *JWT) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.WriteError(w, ErrForbidden, http.StatusForbidden)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(j.secret), nil
		})

		if err != nil {
			response.WriteError(w, ErrInvalidToken, http.StatusForbidden)
			return
		}

		if !token.Valid {
			response.WriteError(w, ErrInvalidToken, http.StatusForbidden)
			return
		}

		userID, ok := claims[UserIDKey].(string)
		if !ok {
			response.WriteError(w, ErrInvalidToken, http.StatusForbidden)
			return
		}

		role, ok := claims[RoleKey].(string)
		if !ok {
			response.WriteError(w, ErrInvalidToken, http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
