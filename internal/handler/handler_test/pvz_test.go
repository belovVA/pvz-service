package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pvz-service/internal/handler"
	"pvz-service/internal/handler/mocks"
	"pvz-service/internal/model"
)

func TestPvzHandler_Create(t *testing.T) {
	mockPvzService := new(mocks.PvzService)
	pvzHandler := handler.NewPvzHandler(mockPvzService)
	r := chi.NewRouter()
	r.Post("/pvz", pvzHandler.CreateNewPvz)
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)

	testPvzID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	tests := []struct {
		name           string
		reqBody        string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "успешное создание",
			reqBody: fmt.Sprintf(`{"city": "%s"}`, handler.KazanRU),
			mockSetup: func() {
				mockPvzService.On("AddNewPvz",
					mock.Anything, model.Pvz{City: handler.KazanRU}).
					Return(&model.Pvz{
						ID:               testPvzID,
						RegistrationDate: fixedTime,
						City:             handler.KazanRU,
					}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: fmt.Sprintf(`{"id":"%s","registrationDate" :"%s","city":"%s"}`,
				testPvzID.String(), fixedTime.Format(time.RFC3339), handler.KazanRU),
		},
		{
			name:           "ошибка создания - invalidBody",
			reqBody:        fmt.Sprintf(`{city: "%s"}`, handler.KazanRU),
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrBodyRequest),
		},
		{
			name:           "ошибка создания - invalidFields",
			reqBody:        fmt.Sprintf(`{"role":"employee"}`),
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrRequestFields),
		},
		{
			name:           "ошибка создания - InvalidCity",
			reqBody:        fmt.Sprintf(`{"city": "%s"}`, "Nizhnekamsk"),
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrInvalidCity),
		},
		{
			name:    "ошибка создания - ошибка сервера",
			reqBody: fmt.Sprintf(`{"city": "%s"}`, handler.SpbRU),
			mockSetup: func() {
				mockPvzService.On("AddNewPvz",
					mock.Anything, model.Pvz{City: handler.SpbRU}).
					Return(nil, fmt.Errorf("server error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"failed to create PVZ: server error"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req, err := http.NewRequest(http.MethodPost, "/pvz", strings.NewReader(tt.reqBody))
			if err != nil {
				t.Fatalf("Ошибка при создании запроса: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			mockPvzService.AssertExpectations(t)
		})
	}
}
