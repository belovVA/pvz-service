package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"pvz-service/internal/model"
)

const (
	PvzOrReceptionsNotExist = "pvz doesn't exists or it doesn't has any receptions"
	ProductNotFound         = "products do not exist"
	FailedProductDelete     = "failed to delete product"
	FailedProductCreate     = "failed to create product"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, typeProduct string, recepID uuid.UUID) (uuid.UUID, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*model.Product, error)
	GetLastProduct(ctx context.Context, receptionID uuid.UUID) (*model.Product, error)
	DeleteProductByID(ctx context.Context, id uuid.UUID) error
	GetProductSliceByReceptionID(ctx context.Context, receptionID uuid.UUID) ([]model.Product, error)
}

type ProductService struct {
	productRepository   ProductRepository
	receptionRepository ReceptionRepository
}

func NewProductService(repoProduct ProductRepository, repoRepository ReceptionRepository) *ProductService {
	return &ProductService{
		productRepository:   repoProduct,
		receptionRepository: repoRepository,
	}
}

func (s *ProductService) AddProduct(ctx context.Context, product model.Product, pvz model.Pvz) (*model.Product, error) {
	reception, err := s.receptionRepository.GetLastReception(ctx, pvz.ID)
	if err != nil {
		return nil, fmt.Errorf(PvzOrReceptionsNotExist)
	}

	if reception.IsClosed {
		return nil, fmt.Errorf(ReceptionAlreadyClosed)
	}

	idProduct, err := s.productRepository.CreateProduct(ctx, product.TypeProduct, reception.ID)
	if err != nil {
		return nil, fmt.Errorf(FailedProductCreate)
	}

	productAns, err := s.productRepository.GetProductByID(ctx, idProduct)
	if err != nil {
		return nil, fmt.Errorf(FailedProductCreate)
	}

	return productAns, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, pvz model.Pvz) error {
	reception, err := s.receptionRepository.GetLastReception(ctx, pvz.ID)
	if err != nil {
		return fmt.Errorf(PvzOrReceptionsNotExist)
	}

	if reception.IsClosed {
		return fmt.Errorf(ReceptionAlreadyClosed)
	}

	product, err := s.productRepository.GetLastProduct(ctx, reception.ID)
	if err != nil {
		return fmt.Errorf(ProductNotFound)
	}

	if err = s.productRepository.DeleteProductByID(ctx, product.ID); err != nil {
		return fmt.Errorf("%s: %s", FailedProductDelete, err.Error())
	}

	return nil
}
