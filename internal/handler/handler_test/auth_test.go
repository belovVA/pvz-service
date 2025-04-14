package handler_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pvz-service/internal/handler"
	"pvz-service/internal/handler/mocks"
	"pvz-service/internal/model"
)

func TestAuthHandler_Register(t *testing.T) {
	mockAuthService := new(mocks.AuthService)
	authHandler := handler.NewAuthHandler(mockAuthService)

	r := chi.NewRouter()
	r.Post("/register", authHandler.Register)

	tests := []struct {
		name           string
		reqBody        string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "успешная регистрация",
			reqBody: `{"email": "test@example.com", "password": "password123", "role": "employee"}`,
			mockSetup: func() {
				mockAuthService.On("Registration", mock.Anything, model.User{
					Email:    "test@example.com",
					Password: "password123",
					Role:     "employee",
				}).Return(&model.User{
					Email:    "test@example.com",
					Password: "password123",
					Role:     "employee",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"00000000-0000-0000-0000-000000000000","email":"test@example.com","role":"employee"}`,
		},
		{
			name:    "ошибка регистрации - invalid role",
			reqBody: `{"email": "test@example.com", "password": "password123", "role": "admin"}`,
			mockSetup: func() {
				mockAuthService.On("Registration", mock.Anything, model.User{
					Email:    "test@example.com",
					Password: "password123",
					Role:     "employee",
				}).Return(nil, errors.New(""))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrInvalidRole),
		},
		{
			name:    "ошибка регистрации - неверные поля запроса",
			reqBody: `{"email": "test@example.com", "role": "employee"}`,
			mockSetup: func() {
				mockAuthService.On("Registration", mock.Anything, model.User{
					Email:    "test@example.com",
					Password: "password123",
					Role:     "employee",
				}).Return(nil, errors.New("registration failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrRequestFields),
		},
		{
			name:    "ошибка регистрации - неверное тело запроса",
			reqBody: `{"email": "test@example.com", :'inv'\rr\t "role": "employee"}`,
			mockSetup: func() {
				mockAuthService.On("Registration", mock.Anything, model.User{
					Email:    "test@example.com",
					Password: "password123",
					Role:     "employee",
				}).Return(nil, errors.New(""))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrBodyRequest),
		},
		{
			name:    "ошибка регистрация - ошибка сервера",
			reqBody: `{"email": "test2@example.com", "password": "password123", "role": "employee"}`,
			mockSetup: func() {
				mockAuthService.On("Registration", mock.Anything, model.User{
					Email:    "test2@example.com",
					Password: "password123",
					Role:     "employee",
				}).Return(nil, fmt.Errorf("server error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req, err := http.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.reqBody))
			if err != nil {
				t.Fatalf("Ошибка при создании запроса: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Authenticate(t *testing.T) {
	mockAuthService := new(mocks.AuthService)
	authHandler := handler.NewAuthHandler(mockAuthService)

	r := chi.NewRouter()
	r.Post("/login", authHandler.Login)

	tests := []struct {
		name           string
		reqBody        string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "успешная авторизация",
			reqBody: `{"email": "test@example.com", "password": "password123"}`,
			mockSetup: func() {
				mockAuthService.On("Authenticate", mock.Anything, model.User{
					Email:    "test@example.com",
					Password: "password123",
				}).Return("some-jwt-token",
					nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `some-jwt-token`,
		},
		{
			name:    "ошибка авторизации - неверное тело запроса",
			reqBody: `{"email": "test@example.com", "role": "employee"}`,
			mockSetup: func() {
				mockAuthService.On("Authenticate", mock.Anything, model.User{
					Email:    "test@example.com",
					Password: "password123",
				}).Return(mock.Anything, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrRequestFields),
		},
		{
			name:    "ошибка авторизации - неверное тело запроса",
			reqBody: `{"email": "test@example.com", :'inv'\rr\t "role": "employee"}`,
			mockSetup: func() {
				mockAuthService.On("Authenticate", mock.Anything, model.User{
					Email:    "test@example.com",
					Password: "password123",
				}).Return(mock.Anything, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrBodyRequest),
		},
		{
			name:    "ошибка авторизациии - ошибка сервера",
			reqBody: `{"email": "test2@example.com", "password": "password123"}`,
			mockSetup: func() {
				mockAuthService.On("Authenticate", mock.Anything, model.User{
					Email:    "test2@example.com",
					Password: "password123",
				}).Return("", fmt.Errorf("server error"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   fmt.Sprintf(`{"message":"server error"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.reqBody))
			if err != nil {
				t.Fatalf("Ошибка при создании запроса: %v", err)
			}
			req.Header.Set("Content-Type", "application/text")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			} else {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}

			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_DummyLogin(t *testing.T) {
	mockAuthService := new(mocks.AuthService)
	authHandler := handler.NewAuthHandler(mockAuthService)

	r := chi.NewRouter()
	r.Post("/dummyLogin", authHandler.DummyLogin)

	tests := []struct {
		name           string
		reqBody        string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "успешная тестовая авторизация",
			reqBody: `{"role": "employee"}`,
			mockSetup: func() {
				mockAuthService.On("DummyAuth", mock.Anything, model.User{Role: handler.EmployeeRole}).Return("some-jwt-token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `some-jwt-token`,
		},
		{
			name:    "ошибка тестовой авторизации - invalid role",
			reqBody: `{"role": "admin"}`,
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrInvalidRole),
		},
		{
			name:    "ошибка тестовой авторизации - неверные поля запроса",
			reqBody: `{}`,
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrRequestFields),
		},
		{
			name:    "ошибка тестовой авторизации - неверное тело запроса",
			reqBody: `{"email": "test@example.com", :'inv'\rr\t "role": "employee"}`,
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrBodyRequest),
		},
		{
			name:    "ошибочная тестовая авторизация - server errpr",
			reqBody: `{"role": "moderator"}`,
			mockSetup: func() {
				mockAuthService.On("DummyAuth", mock.Anything, model.User{Role: handler.ModeratorRole}).Return(
					"", fmt.Errorf("server error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"server error"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req, err := http.NewRequest(http.MethodPost, "/dummyLogin", strings.NewReader(tt.reqBody))
			if err != nil {
				t.Fatalf("Ошибка при создании запроса: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			} else {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}

			mockAuthService.AssertExpectations(t)
		})
	}
}
