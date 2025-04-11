package pvz

import (
	"context"

	"github.com/google/uuid"
	"pvz-service/internal/model"
)

type PvzRepository interface {
	CreatePvz(ctx context.Context, city string) (uuid.UUID, error)
	GetPvzByID(ctx context.Context, id uuid.UUID) (*model.Pvz, error)
}

type PvzService struct {
	pvzRepository PvzRepository
}

func NewPvzService(repo PvzRepository) *PvzService {
	return &PvzService{pvzRepository: repo}
}
