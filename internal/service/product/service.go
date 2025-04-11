package product

import (
	"context"

	"github.com/google/uuid"
	"pvz-service/internal/model"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, typeProduct string, recepID uuid.UUID) (uuid.UUID, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*model.Product, error)
}

type ReceptionRepository interface {
	CreateReception(ctx context.Context, pvzID uuid.UUID) (uuid.UUID, error)
	GetReceptionByID(ctx context.Context, id uuid.UUID) (*model.Reception, error)
	GetLastReception(ctx context.Context, pvzID uuid.UUID) (*model.Reception, error)
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
