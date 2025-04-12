package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"pvz-service/internal/handler/pkg"
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
			pkg.WriteError(w, "Forbidden", http.StatusForbidden)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(j.secret), nil
		})

		if err != nil {
			pkg.WriteError(w, "Invalid token", http.StatusForbidden)
			return
		}

		if !token.Valid {
			pkg.WriteError(w, "Invalid token", http.StatusForbidden)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			pkg.WriteError(w, "Invalid token", http.StatusForbidden)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			pkg.WriteError(w, "Invalid token", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "role", role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
