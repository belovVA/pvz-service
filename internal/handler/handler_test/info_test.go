package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"pvz-service/internal/handler"
	"pvz-service/internal/handler/mocks"
	"pvz-service/internal/model"
)

func TestInfoHandler_GetInfo(t *testing.T) {
	type testCase struct {
		name           string
		queryParams    string
		mockSetup      func(s *mocks.InfoService)
		expectedStatus int
		expectedBody   string
	}

	startDate := "2024-01-01T00:00:00Z"
	endDate := "2025-01-01T00:00:00Z"
	validQuery := "/info?startDate=" + startDate + "&endDate=" + endDate + "&page=1&limit=5"
	invalidDateQuery := "/info?startDate=not-a-date"
	errorQuery := "/info?startDate=" + startDate + "&endDate=" + endDate

	tests := []testCase{
		{
			name:        "успешный запрос",
			queryParams: validQuery,
			mockSetup: func(s *mocks.InfoService) {
				s.On("GetInfoPvz", mock.Anything, &model.PvzInfoQuery{
					StartDate: parseRFC3339(startDate),
					EndDate:   parseRFC3339(endDate),
					Page:      1,
					Limit:     5,
				}).Return([]*model.Pvz{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "невалидная дата",
			queryParams:    invalidDateQuery,
			mockSetup:      func(s *mocks.InfoService) {}, // сервис не вызывается
			expectedStatus: http.StatusBadRequest,
			expectedBody:   handler.ErrConvertParams,
		},
		{
			name:        "ошибка сервиса",
			queryParams: errorQuery,
			mockSetup: func(s *mocks.InfoService) {
				s.On("GetInfoPvz", mock.Anything, &model.PvzInfoQuery{
					StartDate: parseRFC3339(startDate),
					EndDate:   parseRFC3339(endDate),
					Page:      1,
					Limit:     10, // по умолчанию
				}).Return(nil, errors.New("some error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Failed to get PVZ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInfoService := new(mocks.InfoService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockInfoService)
			}

			infoHandler := handler.NewInfoHandler(mockInfoService)

			r := chi.NewRouter()
			r.Get("/info", infoHandler.GetInfo)

			req := httptest.NewRequest(http.MethodGet, tt.queryParams, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}

			mockInfoService.AssertExpectations(t)
		})
	}
}

func parseRFC3339(dateStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, dateStr)
	return t
}
