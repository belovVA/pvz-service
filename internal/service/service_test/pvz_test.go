package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pvz-service/internal/model"
	service2 "pvz-service/internal/service"
	"pvz-service/internal/service/mocks"
)

func TestPvzService_AddNewPvz(t *testing.T) {
	tests := []struct {
		name           string
		city           string
		mockCreatePvz  func(mockRepo *mocks.PvzRepository)
		mockGetPvzByID func(mockRepo *mocks.PvzRepository)
		expectedPvz    *model.Pvz
		expectedError  error
	}{
		{
			name: "Add New Pvz Successfully",
			city: "Moscow",
			mockCreatePvz: func(mockRepo *mocks.PvzRepository) {
				mockRepo.On("CreatePvz", mock.Anything, "Moscow").Return(uuid.New(), nil)
			},
			mockGetPvzByID: func(mockRepo *mocks.PvzRepository) {
				mockRepo.On("GetPvzByID", mock.Anything, mock.Anything).Return(&model.Pvz{
					ID:   uuid.New(),
					City: "Moscow",
				}, nil)
			},
			expectedPvz: &model.Pvz{
				ID:   uuid.New(),
				City: "Moscow",
			},
			expectedError: nil,
		},
		{
			name: "Error Creating Pvz",
			city: "Moscow",
			mockCreatePvz: func(mockRepo *mocks.PvzRepository) {
				mockRepo.On("CreatePvz", mock.Anything, "Moscow").Return(uuid.UUID{}, errors.New("creation failed"))
			},
			mockGetPvzByID: func(mockRepo *mocks.PvzRepository) {},
			expectedPvz:    nil,
			expectedError:  errors.New("creation failed"),
		},
		{
			name: "Error Getting Pvz by ID",
			city: "Moscow",
			mockCreatePvz: func(mockRepo *mocks.PvzRepository) {
				mockRepo.On("CreatePvz", mock.Anything, "Moscow").Return(uuid.New(), nil)
			},
			mockGetPvzByID: func(mockRepo *mocks.PvzRepository) {
				mockRepo.On("GetPvzByID", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))
			},
			expectedPvz:   nil,
			expectedError: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewPvzRepository(t)
			service := service2.NewPvzService(mockRepo)

			// Настроим моки
			tt.mockCreatePvz(mockRepo)
			tt.mockGetPvzByID(mockRepo)

			// Выполняем тестируемую функцию
			pvz, err := service.AddNewPvz(context.Background(), model.Pvz{City: tt.city})

			// Проверяем ошибки
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPvz.City, pvz.City)
			}

			// Проверка, что все ожидания мока были выполнены
			mockRepo.AssertExpectations(t)
		})
	}
}
