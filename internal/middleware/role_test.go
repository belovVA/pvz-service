package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	moderatorRole = "moderator"
	employeeRole  = "employee"
)

func TestRequireRoles(t *testing.T) {
	tests := []struct {
		name           string
		allowedRoles   []string
		userRole       string
		expectedStatus int
	}{
		{
			name:           "user with valid role",
			allowedRoles:   []string{moderatorRole, employeeRole},
			userRole:       employeeRole,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user with invalid role",
			allowedRoles:   []string{employeeRole, employeeRole},
			userRole:       "guest",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "no role",
			allowedRoles:   []string{employeeRole, employeeRole},
			userRole:       "",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "multiple allowed roles",
			allowedRoles:   []string{employeeRole, moderatorRole},
			userRole:       moderatorRole,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "single allowed roles",
			allowedRoles:   []string{moderatorRole},
			userRole:       moderatorRole,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			ctx := context.WithValue(req.Context(), "role", tt.userRole)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler := RequireRoles(tt.allowedRoles...)(next)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
