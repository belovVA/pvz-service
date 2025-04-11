package reception

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"pvz-service/internal/model"
)

func (s *ReceptionService) CloseReception(ctx context.Context, pvzID uuid.UUID) (*model.Reception, error) {
	reception, err := s.receptionRepository.GetLastReception(ctx, pvzID)
	if err != nil {
		return nil, fmt.Errorf("pvz doesn't exist or it has no receptions yet")
	}

	if reception.IsClosed {
		return nil, fmt.Errorf("reception  has  been already closed in this pvz.")
	}

	if err = s.receptionRepository.CloseReception(ctx, reception.ID); err != nil {
		return nil, err
	}

	reception.IsClosed = true

	return reception, nil
}
