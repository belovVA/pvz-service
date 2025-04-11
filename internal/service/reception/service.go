package reception

import (
	"context"

	"github.com/google/uuid"
	"pvz-service/internal/model"
)

type ReceptionRepository interface {
	CreateReception(ctx context.Context, pvzID uuid.UUID) (uuid.UUID, error)
	GetReceptionByID(ctx context.Context, id uuid.UUID) (*model.Reception, error)
	GetLastReception(ctx context.Context, pvzID uuid.UUID) (*model.Reception, error)
	CloseReception(ctx context.Context, receptionID uuid.UUID) error
}

type ReceptionService struct {
	receptionRepository ReceptionRepository
}

func NewReceptionService(repo ReceptionRepository) *ReceptionService {
	return &ReceptionService{
		receptionRepository: repo,
	}
}
