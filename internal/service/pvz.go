package service

import (
	"context"

	"github.com/google/uuid"
	"pvz-service/internal/model"
)

type PvzRepository interface {
	CreatePvz(ctx context.Context, city string) (uuid.UUID, error)
	GetPvzByID(ctx context.Context, id uuid.UUID) (*model.Pvz, error)
	GetIDListPvz(ctx context.Context) ([]uuid.UUID, error)
}

type PvzService struct {
	pvzRepository PvzRepository
}

func NewPvzService(repo PvzRepository) *PvzService {
	return &PvzService{pvzRepository: repo}
}

func (s *PvzService) AddNewPvz(ctx context.Context, pvzModel model.Pvz) (*model.Pvz, error) {
	idPvz, err := s.pvzRepository.CreatePvz(ctx, pvzModel.City)
	if err != nil {
		return nil, err
	}

	pvz, err := s.pvzRepository.GetPvzByID(ctx, idPvz)
	if err != nil {
		return nil, err
	}

	return pvz, nil
}
