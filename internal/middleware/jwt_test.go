package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func mockGenerateToken(t *testing.T, claims map[string]interface{}, secret string) string {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		UserIDKey: claims[UserIDKey],
		RoleKey:   claims[RoleKey],
		"exp":     time.Now().Add(time.Minute).Unix(),
	})

	tokenStr, err := tok.SignedString([]byte(secret))
	assert.NoError(t, err)
	return tokenStr
}

func TestAuthenticate(t *testing.T) {
	secret := "mysecret"
	jwtMiddleware := NewJWT(secret)
	id := uuid.New()
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
				UserIDKey: id.String(),
				RoleKey:   moderatorRole,
			}, secret),
			expectedStatus: http.StatusOK,
			expectedUserID: id.String(),
			expectedRole:   moderatorRole,
		},
		{
			name: "missing userId in token",
			authHeader: "Bearer " + mockGenerateToken(t, map[string]interface{}{
				RoleKey: moderatorRole,
			}, secret),
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "missing role in token",
			authHeader: "Bearer " + mockGenerateToken(t, map[string]interface{}{
				UserIDKey: id.String(),
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
				assert.Equal(t, tt.expectedUserID, r.Context().Value(UserIDKey))
				assert.Equal(t, tt.expectedRole, r.Context().Value(RoleKey))
				w.WriteHeader(http.StatusOK)
			})

			handler := jwtMiddleware.Authenticate(next)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
