package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pvz-service/internal/model"
	service2 "pvz-service/internal/service"
	"pvz-service/internal/service/mocks"
)

func TestInfoService_GetInfoPvz(t *testing.T) {
	// Хардкодим UUID, чтобы был один и тот же во всех структурах
	fixedTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	testPvzID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	testReceptionID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	testProduct1ID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	testProduct2ID := uuid.MustParse("44444444-4444-4444-4444-444444444444")

	tests := []struct {
		name           string
		query          *model.PvzInfoQuery
		mockReceptions []model.Reception
		mockProducts   []model.Product
		mockPvzList    []uuid.UUID
		mockPvzMap     map[uuid.UUID]*model.Pvz
		expectedResult []*model.Pvz
		expectedError  error
	}{
		{
			name: "success",
			query: &model.PvzInfoQuery{
				StartDate: time.Time{},
				EndDate:   time.Time{},
				Page:      1,
				Limit:     10,
			},
			mockReceptions: []model.Reception{
				{ID: testReceptionID, DateTime: fixedTime, IsClosed: false, PvzID: testPvzID},
			},
			mockProducts: []model.Product{
				{ID: testProduct1ID, TypeProduct: "Product1", ReceptionID: testReceptionID},
				{ID: testProduct2ID, TypeProduct: "Product2", ReceptionID: testReceptionID},
			},
			mockPvzList: []uuid.UUID{
				testPvzID,
			},
			mockPvzMap: map[uuid.UUID]*model.Pvz{
				testPvzID: {
					ID:   testPvzID,
					City: "City1",
					Receptions: []model.Reception{
						{ID: testReceptionID, DateTime: fixedTime, PvzID: testPvzID, IsClosed: false},
					},
				},
			},
			expectedResult: []*model.Pvz{
				{
					ID:   testPvzID,
					City: "City1",
					Receptions: []model.Reception{
						{
							ID:       testReceptionID,
							DateTime: fixedTime,
							PvzID:    testPvzID,
							IsClosed: false,
							Products: []model.Product{
								{
									ID:          testProduct1ID,
									TypeProduct: "Product1",
									ReceptionID: testReceptionID,
								},
								{
									ID:          testProduct2ID,
									TypeProduct: "Product2",
									ReceptionID: testReceptionID,
								},
							},
						},
					},
				},
			},

			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReceptionRepo := new(mocks.ReceptionRepository)
			mockProductRepo := new(mocks.ProductRepository)
			mockPvzRepo := new(mocks.PvzRepository)

			mockReceptionRepo.On("GetReceptionsSliceWithTimeRange", mock.Anything, tt.query.StartDate, tt.query.EndDate).Return(tt.mockReceptions, nil)
			mockProductRepo.On("GetProductSliceByReceptionID", mock.Anything, mock.Anything).Return(tt.mockProducts, nil)

			mockPvzRepo.On("GetIDListPvz", mock.Anything).Return(tt.mockPvzList, nil)
			mockPvzRepo.On("GetPvzByID", mock.Anything, testPvzID).Return(tt.mockPvzMap[testPvzID], nil)

			service := service2.NewInfoService(mockProductRepo, mockReceptionRepo, mockPvzRepo)

			result, err := service.GetInfoPvz(context.Background(), tt.query)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)

			mockReceptionRepo.AssertExpectations(t)
			mockProductRepo.AssertExpectations(t)
			mockPvzRepo.AssertExpectations(t)
		})
	}
}
