package reception

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"pvz-service/internal/model"
)

func (s *ReceptionService) CreateReception(ctx context.Context, pvzID uuid.UUID) (*model.Reception, error) {
	// Проверяем наличие последней приемки в данном ПВЗ и смотрим, был ли он закрыт
	reception, err := s.receptionRepository.GetLastReception(ctx, pvzID)
	if err == nil {
		if !reception.IsClosed {
			return nil, fmt.Errorf("in this PVZ, the reception  has not been closed yet.")
		}
	}

	id, err := s.receptionRepository.CreateReception(ctx, pvzID)
	if err != nil {
		return nil, err
	}

	rep, err := s.receptionRepository.GetReceptionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return rep, nil
}
