package handler_test

import (
	"bytes"
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

func TestProductHandlers_CreateNewProduct(t *testing.T) {
	mockService := new(mocks.ProductService)
	handl := handler.NewProductHandler(mockService)

	router := chi.NewRouter()
	router.Post("/product", handl.CreateNewProduct)

	pvzID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	testRecepID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	product := &model.Product{
		ID:          uuid.New(),
		ReceptionID: testRecepID,
		TypeProduct: handler.ElectrType,
		DateTime:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name           string
		body           string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "успешное создание",
			body: fmt.Sprintf(`{"type":"%s","pvzId":"%s"}`, handler.ElectrType, pvzID),
			mockSetup: func() {
				mockService.On("AddProduct", mock.Anything, handler.ElectrType, pvzID).Return(product, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   fmt.Sprintf(`{"id":"%s","receptionId":"%s","type":"%s","dateTime":"%s"}`, product.ID, testRecepID, handler.ElectrType, product.DateTime.Format(time.RFC3339)),
		},
		{
			name: "невалидное тело запроса",
			body: `{typeProduct:"electronic"}`, // неправильно оформлен JSON
			mockSetup: func() {
				// нет вызова сервиса
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrBodyRequest),
		},
		{
			name: "невалидный UUID",
			body: `{"type":"electronic","pvzId":"not-a-uuid"}`,
			mockSetup: func() {
				// нет вызова сервиса
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrUUIDParsing),
		},
		{
			name:           "неподдерживаемый тип продукта",
			body:           fmt.Sprintf(`{"type":"weird","pvzId":"%s"}`, pvzID),
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrProductType),
		},
		{
			name: "ошибка при создании продукта",
			body: fmt.Sprintf(`{"type":"%s","pvzId":"%s"}`, handler.ClothesType, pvzID),
			mockSetup: func() {
				mockService.On("AddProduct", mock.Anything, handler.ClothesType, pvzID).Return(nil, errors.New("DB error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Failed add Product: DB error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestProductHandlers_RemoveLastProduct(t *testing.T) {
	mockService := new(mocks.ProductService)
	handl := handler.NewProductHandler(mockService)

	router := chi.NewRouter()
	router.Delete("/product/{pvzId}", handl.RemoveLastProduct)

	validPvzID := uuid.New()
	validPvzID2 := uuid.New()

	tests := []struct {
		name           string
		pvzIdPath      string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "успешное удаление",
			pvzIdPath: validPvzID.String(),
			mockSetup: func() {
				mockService.On("DeleteProduct", mock.Anything, validPvzID).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   ``,
		},
		{
			name:           "невалидный UUID",
			pvzIdPath:      "invalid-uuid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   fmt.Sprintf(`{"message":"%s"}`, handler.ErrUUIDParsing),
		},
		{
			name:      "ошибка при удалении",
			pvzIdPath: validPvzID2.String(),
			mockSetup: func() {
				mockService.On("DeleteProduct", mock.Anything, validPvzID2).Return(errors.New("delete error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"failed to delete product:  delete error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/product/%s", tt.pvzIdPath), nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedBody, w.Body.String())
			} else {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())

			}

			mockService.AssertExpectations(t)
		})
	}
}
