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

//const (
//	PvzOrReceptionsNotExist = "pvz doesn't exists or it doesn't has any receptions"
//	ProductNotFound         = "products do not exist"
//	FailedProductDelete     = "failed to delete product"
//	ReceptionAlreadyClosed  = "reception has been already closed in this pvz"

//)

const electrType = "электроника"

func TestProductService_AddProduct(t *testing.T) {
	tests := []struct {
		name                 string
		pvzID                uuid.UUID
		typeProduct          string
		mockGetLastReception func(mockRepo *mocks.ReceptionRepository)
		mockCreateProduct    func(mockRepo *mocks.ProductRepository)
		mockGetProductByID   func(mockRepo *mocks.ProductRepository)
		expectedError        error
		expectedProduct      *model.Product
	}{
		{
			name:        "Add Product when last reception is closed",
			pvzID:       uuid.New(),
			typeProduct: electrType,
			mockGetLastReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetLastReception", mock.Anything, mock.Anything).Return(&model.Reception{IsClosed: true}, nil)
			},
			mockCreateProduct:  func(mockRepo *mocks.ProductRepository) {},
			mockGetProductByID: func(mockRepo *mocks.ProductRepository) {},
			expectedError:      errors.New(service.ReceptionAlreadyClosed),
			expectedProduct:    nil,
		},
		{
			name:        "Add Product successfully",
			pvzID:       uuid.New(),
			typeProduct: electrType,
			mockGetLastReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetLastReception", mock.Anything, mock.Anything).Return(&model.Reception{IsClosed: false, ID: uuid.New()}, nil)
			},
			mockCreateProduct: func(mockRepo *mocks.ProductRepository) {
				mockRepo.On("CreateProduct", mock.Anything, mock.Anything, mock.Anything).Return(uuid.New(), nil)
			},
			mockGetProductByID: func(mockRepo *mocks.ProductRepository) {
				mockRepo.On("GetProductByID", mock.Anything, mock.Anything).Return(&model.Product{ID: uuid.New(), TypeProduct: electrType}, nil)
			},
			expectedError:   nil,
			expectedProduct: &model.Product{ID: uuid.New(), TypeProduct: electrType},
		},
		{
			name:        "Error when reception not found",
			pvzID:       uuid.New(),
			typeProduct: electrType,
			mockGetLastReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetLastReception", mock.Anything, mock.Anything).Return(nil, errors.New(service.PvzOrReceptionsNotExist))
			},
			mockCreateProduct:  func(mockRepo *mocks.ProductRepository) {},
			mockGetProductByID: func(mockRepo *mocks.ProductRepository) {},
			expectedError:      errors.New(service.PvzOrReceptionsNotExist),
			expectedProduct:    nil,
		},
		{
			name:        "Error when product creation fails",
			pvzID:       uuid.New(),
			typeProduct: electrType,
			mockGetLastReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetLastReception", mock.Anything, mock.Anything).Return(&model.Reception{IsClosed: false, ID: uuid.New()}, nil)
			},
			mockCreateProduct: func(mockRepo *mocks.ProductRepository) {
				mockRepo.On("CreateProduct", mock.Anything, mock.Anything, mock.Anything).Return(uuid.UUID{}, errors.New("product creation failed"))
			},
			mockGetProductByID: func(mockRepo *mocks.ProductRepository) {},
			expectedError:      errors.New("product creation failed"),
			expectedProduct:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProductRepo := mocks.NewProductRepository(t)
			mockReceptionRepo := mocks.NewReceptionRepository(t)
			service := service.NewProductService(mockProductRepo, mockReceptionRepo)

			// Настроим моки
			tt.mockGetLastReception(mockReceptionRepo)
			tt.mockCreateProduct(mockProductRepo)
			tt.mockGetProductByID(mockProductRepo)

			// Выполняем тестируемую функцию
			product, err := service.AddProduct(context.Background(), tt.typeProduct, tt.pvzID)

			// Проверяем ошибки
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedProduct.TypeProduct, product.TypeProduct)
			}

			// Проверка, что все ожидания мока были выполнены
			mockProductRepo.AssertExpectations(t)
			mockReceptionRepo.AssertExpectations(t)
		})
	}
}

func TestProductService_DeleteProduct(t *testing.T) {
	tests := []struct {
		name                  string
		pvzID                 uuid.UUID
		mockGetLastReception  func(mockRepo *mocks.ReceptionRepository)
		mockGetLastProduct    func(mockRepo *mocks.ProductRepository)
		mockDeleteProductByID func(mockRepo *mocks.ProductRepository)
		expectedError         error
	}{
		{
			name:  "Delete Product when no reception found",
			pvzID: uuid.New(),
			mockGetLastReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetLastReception", mock.Anything, mock.Anything).Return(nil, errors.New(service.PvzOrReceptionsNotExist))
			},
			mockGetLastProduct:    func(mockRepo *mocks.ProductRepository) {},
			mockDeleteProductByID: func(mockRepo *mocks.ProductRepository) {},
			expectedError:         errors.New(service.PvzOrReceptionsNotExist),
		},
		{
			name:  "Delete Product when reception is closed",
			pvzID: uuid.New(),
			mockGetLastReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetLastReception", mock.Anything, mock.Anything).Return(&model.Reception{IsClosed: true}, nil)
			},
			mockGetLastProduct:    func(mockRepo *mocks.ProductRepository) {},
			mockDeleteProductByID: func(mockRepo *mocks.ProductRepository) {},
			expectedError:         errors.New(service.ReceptionAlreadyClosed),
		},
		{
			name:  "Delete Product successfully",
			pvzID: uuid.New(),
			mockGetLastReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetLastReception", mock.Anything, mock.Anything).Return(&model.Reception{IsClosed: false, ID: uuid.New()}, nil)
			},
			mockGetLastProduct: func(mockRepo *mocks.ProductRepository) {
				mockRepo.On("GetLastProduct", mock.Anything, mock.Anything).Return(&model.Product{ID: uuid.New()}, nil)
			},
			mockDeleteProductByID: func(mockRepo *mocks.ProductRepository) {
				mockRepo.On("DeleteProductByID", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "Error when product not found",
			pvzID: uuid.New(),
			mockGetLastReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetLastReception", mock.Anything, mock.Anything).Return(&model.Reception{IsClosed: false, ID: uuid.New()}, nil)
			},
			mockGetLastProduct: func(mockRepo *mocks.ProductRepository) {
				mockRepo.On("GetLastProduct", mock.Anything, mock.Anything).Return(nil, errors.New(service.ProductNotFound))
			},
			mockDeleteProductByID: func(mockRepo *mocks.ProductRepository) {},
			expectedError:         errors.New(service.ProductNotFound),
		},
		{
			name:  "Error when deleting product fails",
			pvzID: uuid.New(),
			mockGetLastReception: func(mockRepo *mocks.ReceptionRepository) {
				mockRepo.On("GetLastReception", mock.Anything, mock.Anything).Return(&model.Reception{IsClosed: false, ID: uuid.New()}, nil)
			},
			mockGetLastProduct: func(mockRepo *mocks.ProductRepository) {
				mockRepo.On("GetLastProduct", mock.Anything, mock.Anything).Return(&model.Product{ID: uuid.New()}, nil)
			},
			mockDeleteProductByID: func(mockRepo *mocks.ProductRepository) {
				mockRepo.On("DeleteProductByID", mock.Anything, mock.Anything).Return(errors.New("delete failed"))
			},
			expectedError: errors.New(service.FailedProductDelete + ": delete failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProductRepo := mocks.NewProductRepository(t)
			mockReceptionRepo := mocks.NewReceptionRepository(t)
			service := service.NewProductService(mockProductRepo, mockReceptionRepo)

			// Настроим моки
			tt.mockGetLastReception(mockReceptionRepo)
			tt.mockGetLastProduct(mockProductRepo)
			tt.mockDeleteProductByID(mockProductRepo)

			// Выполняем тестируемую функцию
			err := service.DeleteProduct(context.Background(), tt.pvzID)

			// Проверяем ошибки
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Проверка, что все ожидания мока были выполнены
			mockProductRepo.AssertExpectations(t)
			mockReceptionRepo.AssertExpectations(t)
		})
	}
}
