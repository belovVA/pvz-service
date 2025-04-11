package product

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"pvz-service/internal/model"
)

func (s *ProductService) AddProduct(ctx context.Context, typeProduct string, pvzID uuid.UUID) (*model.Product, error) {
	if err := s.validateType(typeProduct); err != nil {
		return nil, err
	}

	reception, err := s.receptionRepository.GetLastReception(ctx, pvzID)
	if err != nil {
		return nil, err
	}

	if reception.IsClosed {
		return nil, fmt.Errorf("in this PVZ, the reception  has been closed already.")
	}

	idProduct, err := s.productRepository.CreateProduct(ctx, typeProduct, reception.ID)
	if err != nil {
		log.Println("qwer")

		return nil, err
	}

	product, err := s.productRepository.GetProductByID(ctx, idProduct)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) validateType(typeProduct string) error {
	switch typeProduct {
	case "электроника":
		return nil
	case "одежда":
		return nil
	case "обувь":
		return nil
	}
	return fmt.Errorf("invalid product type: %s", typeProduct)
}
