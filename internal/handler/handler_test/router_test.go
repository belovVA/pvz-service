package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pvz-service/internal/handler"
	"pvz-service/internal/handler/mocks"
	"pvz-service/pkg/jwtutils"
	"pvz-service/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessControl_AllRoutes(t *testing.T) {
	mockService := new(mocks.Service)
	secret := "test-secret"
	logger := logger.InitLogger()

	r := handler.NewRouter(mockService, secret, logger)

	type testCase struct {
		name           string
		method         string
		path           string
		role           string // "" = без токена
		expectedStatus int
	}

	tests := []testCase{
		// No JWT
		{"NoToken /pvz POST", http.MethodPost, "/pvz", "", http.StatusForbidden},
		{"NoToken /pvz GET", http.MethodGet, "/pvz", "", http.StatusForbidden},
		{"NoToken /receptions POST", http.MethodPost, "/receptions", "", http.StatusForbidden},
		{"NoToken /products POST", http.MethodPost, "/products", "", http.StatusForbidden},
		{"NoToken /pvz/{id}/close_last_reception", http.MethodPost, "/pvz/123/close_last_reception", "", http.StatusForbidden},
		{"NoToken /pvz/{id}/delete_last_product", http.MethodPost, "/pvz/123/delete_last_product", "", http.StatusForbidden},

		//Wrong Role
		{"WrongRole-Employee /pvz POST", http.MethodPost, "/pvz", handler.EmployeeRole, http.StatusForbidden},
		{"InvalidRole /pvz GET", http.MethodGet, "/pvz", "invalid", http.StatusForbidden},
		{"WrongRole-Moderator /receptions POST", http.MethodPost, "/receptions", handler.ModeratorRole, http.StatusForbidden},
		{"WrongRole-Moderator /products POST", http.MethodPost, "/products", handler.ModeratorRole, http.StatusForbidden},
		{"WrongRole-Moderator /pvz/{id}/close_last_reception", http.MethodPost, "/pvz/123/close_last_reception", handler.ModeratorRole, http.StatusForbidden},
		{"WrongRole-Moderator /pvz/{id}/delete_last_product", http.MethodPost, "/pvz/123/delete_last_product", handler.ModeratorRole, http.StatusForbidden},

		// Good Role
		//{"Employee /receptions POST", http.MethodPost, "/receptions", handler.EmployeeRole, http.StatusBadRequest},

		//{"Employee /products POST", http.MethodPost, "/products", handler.EmployeeRole, http.StatusBadRequest},

		//{"Employee /pvz/{id}/close_last_reception", http.MethodPost, "/pvz/123/close_last_reception", handler.EmployeeRole, http.StatusBadRequest},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			req = httptest.NewRequest(tt.method, tt.path, nil)

			if tt.role != "" {

				token, err := jwtutils.Generate(map[string]interface{}{
					"userId": "test-user",
					"role":   tt.role,
				}, time.Hour, secret)
				require.NoError(t, err)
				req.Header.Set("Authorization", "Bearer "+token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
