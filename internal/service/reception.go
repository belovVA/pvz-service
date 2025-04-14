package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"pvz-service/internal/model"
)

const ReceptionWasNotClosed = "in this PVZ, the reception  has not been closed yet."
const ReceptionAlreadyClosed = "reception  has  been already closed in this pvz."

type ReceptionRepository interface {
	CreateReception(ctx context.Context, pvzID uuid.UUID) (uuid.UUID, error)
	GetReceptionByID(ctx context.Context, id uuid.UUID) (*model.Reception, error)
	GetLastReception(ctx context.Context, pvzID uuid.UUID) (*model.Reception, error)
	CloseReception(ctx context.Context, receptionID uuid.UUID) error
	GetReceptionsSliceWithTimeRange(ctx context.Context, begin time.Time, end time.Time) ([]model.Reception, error)
}

type ReceptionService struct {
	receptionRepository ReceptionRepository
}

func NewReceptionService(repo ReceptionRepository) *ReceptionService {
	return &ReceptionService{
		receptionRepository: repo,
	}
}

func (s *ReceptionService) CreateReception(ctx context.Context, receptionModel model.Reception) (*model.Reception, error) {
	// Проверяем наличие последней приемки в данном ПВЗ и смотрим, был ли он закрыт
	reception, err := s.receptionRepository.GetLastReception(ctx, receptionModel.PvzID)

	if err == nil && !reception.IsClosed {
		return nil, fmt.Errorf(ReceptionWasNotClosed)
	}

	// Если ПВЗ с таким ID нет, то Constraint вернет ошибку и приемка не будет создана
	id, err := s.receptionRepository.CreateReception(ctx, receptionModel.PvzID)
	if err != nil {
		return nil, err
	}

	rep, err := s.receptionRepository.GetReceptionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return rep, nil
}

func (s *ReceptionService) CloseReception(ctx context.Context, receptionModel model.Reception) (*model.Reception, error) {
	reception, err := s.receptionRepository.GetLastReception(ctx, receptionModel.PvzID)
	if err != nil {
		return nil, fmt.Errorf(PvzOrReceptionsNotExist)
	}

	if reception.IsClosed {
		return nil, fmt.Errorf(ReceptionAlreadyClosed)
	}

	if err = s.receptionRepository.CloseReception(ctx, reception.ID); err != nil {
		return nil, err
	}

	reception.IsClosed = true

	return reception, nil
}
