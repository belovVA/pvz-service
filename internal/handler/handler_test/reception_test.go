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

func TestReceptionHandlers_OpenNewReception(t *testing.T) {
	mockReceptionService := new(mocks.ReceptionService)
	receptionHandler := handler.NewReceptionHandler(mockReceptionService)
	r := chi.NewRouter()
	r.Post("/receptions", receptionHandler.OpenNewReception)

	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	testPvzID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	testPvzID2 := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	testRecepID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	tests := []struct {
		name           string
		reqBody        string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "успешное создание нового приёма",
			reqBody: fmt.Sprintf(`{"pvzID": "%s"}`, testPvzID),
			mockSetup: func() {
				mockReceptionService.On("CreateReception", mock.Anything, testPvzID).Return(&model.Reception{
					ID:       testRecepID,
					DateTime: fixedTime,
					IsClosed: false,
					PvzID:    testPvzID,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: fmt.Sprintf(`{"id":"%s","dateTime":"%s",  "pvzId":"%s", "status":"%s"}`,
				testRecepID.String(), fixedTime.Format(time.RFC3339), testPvzID.String(), "in_progress"),
		},
		{
			name:    "ошибка при обработке тела запроса - неверный формат",
			reqBody: fmt.Sprintf(`{pvzID: "%s"}`, testPvzID), // неверный формат JSON
			mockSetup: func() {
				mockReceptionService.On("CreateReception", mock.Anything, testPvzID).Return(nil, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrBodyRequest),
		},
		{
			name:    "ошибка при проверке полей запроса - неверные данные",
			reqBody: fmt.Sprintf(`{"role":"employee"}`),
			mockSetup: func() {
				mockReceptionService.On("CreateReception", mock.Anything, testPvzID).Return(nil, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrRequestFields),
		},
		{
			name:    "ошибка при создании - сервис не смог создать",
			reqBody: fmt.Sprintf(`{"pvzID": "%s"}`, testPvzID2),
			mockSetup: func() {
				mockReceptionService.On("CreateReception", mock.Anything, testPvzID2).Return(nil, fmt.Errorf("failed to create reception"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"failed to create Reception: failed to create reception"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req, err := http.NewRequest(http.MethodPost, "/receptions", strings.NewReader(tt.reqBody))
			if err != nil {
				t.Fatalf("Ошибка при создании запроса: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			mockReceptionService.AssertExpectations(t)
		})
	}
}

func TestReceptionHandlers_CloseLastReception(t *testing.T) {
	mockReceptionService := new(mocks.ReceptionService)
	receptionHandler := handler.NewReceptionHandler(mockReceptionService)
	r := chi.NewRouter()

	r.Post("/receptions/{pvzId}/close", receptionHandler.CloseLastReception)

	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	testPvzID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	testPvzID2 := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	testRecepID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	tests := []struct {
		name           string
		pvzIdParam     string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:       "успешное закрытие приёма",
			pvzIdParam: testPvzID.String(),
			mockSetup: func() {
				mockReceptionService.On("CloseReception", mock.Anything, testPvzID).Return(&model.Reception{
					ID:       testRecepID,
					DateTime: fixedTime,
					IsClosed: true,
					PvzID:    testPvzID,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: fmt.Sprintf(`{"id":"%s","dateTime":"%s",  "pvzId":"%s", "status":"%s"}`,
				testRecepID.String(), fixedTime.Format(time.RFC3339), testPvzID.String(), "close"),
		},
		{
			name:       "ошибка при неверном формате ID",
			pvzIdParam: "invalid-uuid",
			mockSetup: func() {
				mockReceptionService.On("CloseReception", mock.Anything, testPvzID).Return(nil, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrUUIDParsing),
		},
		{
			name:       "ошибка при закрытии приёма - сервис не смог закрыть",
			pvzIdParam: testPvzID2.String(),
			mockSetup: func() {
				mockReceptionService.On("CloseReception", mock.Anything, testPvzID2).Return(nil, fmt.Errorf("failed to close reception"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"failed to close reception: failed to close reception"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/receptions/%s/close", tt.pvzIdParam), nil)
			if err != nil {
				t.Fatalf("Ошибка при создании запроса: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			mockReceptionService.AssertExpectations(t)
		})
	}
}
