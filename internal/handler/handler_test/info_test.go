package handler_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestInfoHandler_GetInfo(t *testing.T) {
	type testCase struct {
		name           string
		queryParams    string
		mockSetup      func(s *mocks.InfoService)
		expectedStatus int
		expectedBody   string
	}
	pvz1ID := uuid.New()
	pvz2ID := uuid.New()
	product1ID := uuid.New()
	product2ID := uuid.New()
	product3ID := uuid.New()
	recepID1 := uuid.New()
	recepID2 := uuid.New()
	dateTime := time.Now()
	startDate := "2024-01-01T00:00:00Z"
	endDate := "2025-01-01T00:00:00Z"
	validQuery := "/info?startDate=" + startDate + "&endDate=" + endDate + "&page=1&limit=5"
	invalidDateQuery := "/info?startDate=not-a-date"
	invalidDateQuery2 := "/info?endDate=not-a-date"
	invalidPageQuery := "/info?page=ten"
	errorQuery := "/info?startDate=" + startDate + "&endDate=" + endDate

	product1 := model.Product{
		ID:          product1ID,
		DateTime:    dateTime,
		TypeProduct: handler.ElectrType,
		ReceptionID: recepID1,
	}
	product2 := model.Product{
		ID:          product2ID,
		DateTime:    dateTime,
		TypeProduct: handler.ElectrType,
		ReceptionID: recepID1,
	}
	product3 := model.Product{
		ID:          product3ID,
		DateTime:    dateTime,
		TypeProduct: handler.ElectrType,
		ReceptionID: recepID2,
	}
	reception1 := model.Reception{
		ID:       recepID1,
		DateTime: dateTime,
		Products: []model.Product{product1, product2},
		IsClosed: true,
		PvzID:    pvz1ID,
	}
	reception2 := model.Reception{
		ID:       recepID2,
		DateTime: dateTime,
		Products: []model.Product{product3},
		IsClosed: false,
		PvzID:    pvz1ID,
	}
	pvz := model.Pvz{
		ID:               pvz1ID,
		RegistrationDate: parseRFC3339(startDate),
		City:             "Москва",
		Receptions:       []model.Reception{reception1, reception2},
	}
	pvz2 := model.Pvz{
		ID:               pvz2ID,
		RegistrationDate: parseRFC3339(startDate),
		City:             "Москва",
		Receptions:       []model.Reception{},
	}

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
			mockSetup:      func(s *mocks.InfoService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   handler.ErrConvertParams,
		},

		{
			name:           "невалидная дата2",
			queryParams:    invalidDateQuery2,
			mockSetup:      func(s *mocks.InfoService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   handler.ErrConvertParams,
		},
		{
			name:           "невалидная страница",
			queryParams:    invalidPageQuery,
			mockSetup:      func(s *mocks.InfoService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   handler.ErrQueryParameters,
		},
		{
			name:        "ошибка сервиса",
			queryParams: errorQuery,
			mockSetup: func(s *mocks.InfoService) {
				s.On("GetInfoPvz",
					mock.Anything, &model.PvzInfoQuery{
						StartDate: parseRFC3339(startDate),
						EndDate:   parseRFC3339(endDate),
						Page:      1,
						Limit:     10, // по умолчанию
					}).Return(nil, errors.New("some error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Failed to get PVZ: some error",
		},
		{
			name:        "успешный запрос c несколькими пвз",
			queryParams: validQuery,
			mockSetup: func(s *mocks.InfoService) {
				s.On("GetInfoPvz",
					mock.Anything, &model.PvzInfoQuery{
						StartDate: parseRFC3339(startDate),
						EndDate:   parseRFC3339(endDate),
						Page:      1,
						Limit:     5,
					}).Return([]*model.Pvz{&pvz, &pvz2}, nil)
			},
			expectedStatus: http.StatusOK,
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
				assert.Contains(t, rec.Body.String(),
					fmt.Sprintf(`{"message":"%s"}`, tt.expectedBody))
			}

			mockInfoService.AssertExpectations(t)
		})
	}
}

func parseRFC3339(dateStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, dateStr)
	return t
}
