package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pvz-service/internal/model"
	"pvz-service/internal/service"
	"pvz-service/internal/service/mocks"
)

func TestReceptionService_CreateReception(t *testing.T) {
	tests := []struct {
		name                 string
		pvzID                uuid.UUID
		mockGetLastReception func(mockRepo *mocks.ReceptionRepository)
		mockCreateReception  func(mockRepo *mocks.ReceptionRepository)
		mockGetReceptionByID func(mockRepo *mocks.ReceptionRepository)
		expectedError        error
		expectedReception    *model.Reception
	}{
		{
			name:  "Create Reception when previous one is not closed",
			pvzID: uuid.New(),
			mockGetLastReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetLastReception", mock.Anything, mock.Anything).Return(&model.Reception{IsClosed: false}, nil)
			},
			mockCreateReception:  func(mockRepo *mocks.ReceptionRepository) {},
			mockGetReceptionByID: func(mockRepo *mocks.ReceptionRepository) {},
			expectedError:        errors.New(service.ReceptionWasNotClosed),
			expectedReception:    nil,
		},
		{
			name:  "Create Reception when previous one is closed",
			pvzID: uuid.New(),
			mockGetLastReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetLastReception", mock.Anything, mock.Anything).Return(&model.Reception{IsClosed: true}, nil)
			},
			mockCreateReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("CreateReception", mock.Anything, mock.Anything).Return(uuid.New(), nil)
			},
			mockGetReceptionByID: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetReceptionByID", mock.Anything, mock.Anything).Return(&model.Reception{ID: uuid.New(), IsClosed: false}, nil)
			},
			expectedError:     nil,
			expectedReception: &model.Reception{ID: uuid.New(), IsClosed: false},
		},
		{
			name:  "Error creating reception",
			pvzID: uuid.New(),
			mockGetLastReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetLastReception", mock.Anything, mock.Anything).Return(&model.Reception{IsClosed: true}, nil)
			},
			mockCreateReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("CreateReception", mock.Anything, mock.Anything).Return(uuid.UUID{}, errors.New("creation failed"))
			},
			mockGetReceptionByID: func(mockRepo *mocks.ReceptionRepository) {},
			expectedError:        errors.New("creation failed"),
			expectedReception:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewReceptionRepository(t)
			service := service.NewReceptionService(mockRepo)

			// Настроим моки
			tt.mockGetLastReception(mockRepo)
			tt.mockCreateReception(mockRepo)
			tt.mockGetReceptionByID(mockRepo)

			// Выполняем тестируемую функцию
			reception, err := service.CreateReception(context.Background(), model.Reception{PvzID: tt.pvzID})

			// Проверяем ошибки
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedReception.IsClosed, reception.IsClosed)
			}

			// Проверка, что все ожидания мока были выполнены
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestReceptionService_CloseReception(t *testing.T) {
	tests := []struct {
		name                 string
		pvzID                uuid.UUID
		mockGetLastReception func(mockRepo *mocks.ReceptionRepository)
		mockCloseReception   func(mockRepo *mocks.ReceptionRepository)
		expectedError        error
		expectedReception    *model.Reception
	}{
		{
			name:  "Close Reception when last reception is closed",
			pvzID: uuid.New(),
			mockGetLastReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetLastReception", mock.Anything, mock.Anything).Return(&model.Reception{IsClosed: true}, nil)
			},
			mockCloseReception: func(mockRepo *mocks.ReceptionRepository) {},
			expectedError:      errors.New(service.ReceptionAlreadyClosed),
			expectedReception:  nil,
		},
		{
			name:  "Close Reception successfully",
			pvzID: uuid.New(),
			mockGetLastReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetLastReception", mock.Anything, mock.Anything).Return(&model.Reception{IsClosed: false, ID: uuid.New()}, nil)
			},
			mockCloseReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("CloseReception", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError:     nil,
			expectedReception: &model.Reception{ID: uuid.New(), IsClosed: true},
		},
		{
			name:  "Error when no reception found",
			pvzID: uuid.New(),
			mockGetLastReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetLastReception", mock.Anything, mock.Anything).Return(nil, errors.New(service.PvzOrReceptionsNotExist))
			},
			mockCloseReception: func(mockRepo *mocks.ReceptionRepository) {},
			expectedError:      errors.New(service.PvzOrReceptionsNotExist),
			expectedReception:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewReceptionRepository(t)
			service := service.NewReceptionService(mockRepo)

			// Настроим моки
			tt.mockGetLastReception(mockRepo)
			tt.mockCloseReception(mockRepo)

			// Выполняем тестируемую функцию
			reception, err := service.CloseReception(context.Background(), model.Reception{PvzID: tt.pvzID})

			// Проверяем ошибки
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedReception.IsClosed, reception.IsClosed)
			}

			// Проверка, что все ожидания мока были выполнены
			mockRepo.AssertExpectations(t)
		})
	}
}
