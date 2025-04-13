package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func mockGenerateToken(t *testing.T, claims map[string]interface{}, secret string) string {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims["user_id"],
		"role":    claims["role"],
		"exp":     time.Now().Add(time.Minute).Unix(),
	})

	tokenStr, err := tok.SignedString([]byte(secret))
	assert.NoError(t, err)
	return tokenStr
}

func TestAuthenticate(t *testing.T) {
	secret := "mysecret"
	jwtMiddleware := NewJWT(secret)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedUserID string
		expectedRole   string
	}{
		{
			name:           "missing Authorization header",
			authHeader:     "",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "invalid token format",
			authHeader:     "Bearer invalid.token.string",
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "valid token",
			authHeader: "Bearer " + mockGenerateToken(t, map[string]interface{}{
				"user_id": "123",
				"role":    "admin",
			}, secret),
			expectedStatus: http.StatusOK,
			expectedUserID: "123",
			expectedRole:   "admin",
		},
		{
			name: "missing user_id in token",
			authHeader: "Bearer " + mockGenerateToken(t, map[string]interface{}{
				"role": "admin",
			}, secret),
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "missing role in token",
			authHeader: "Bearer " + mockGenerateToken(t, map[string]interface{}{
				"user_id": "123",
			}, secret),
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.expectedUserID, r.Context().Value("user_id"))
				assert.Equal(t, tt.expectedRole, r.Context().Value("role"))
				w.WriteHeader(http.StatusOK)
			})

			handler := jwtMiddleware.Authenticate(next)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
