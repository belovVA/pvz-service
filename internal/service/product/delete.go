package product

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *ProductService) DeleteProduct(ctx context.Context, pvzID uuid.UUID) error {
	reception, err := s.receptionRepository.GetLastReception(ctx, pvzID)
	if err != nil {
		return fmt.Errorf("pvz doesn't exists or it doesn't has any receptions")
	}

	if reception.IsClosed {
		return fmt.Errorf("in this PVZ, the reception  has been closed already.")
	}

	product, err := s.productRepository.GetLastProduct(ctx, reception.ID)
	if err != nil {
		return fmt.Errorf("products do not exist")
	}

	if err := s.productRepository.DeleteProductByID(ctx, product.ID); err != nil {
		return fmt.Errorf("failed to delete product: %s", err.Error())
	}

	return nil
}
