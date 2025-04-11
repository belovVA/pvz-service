package pvz

import (
	"context"
	"fmt"

	"pvz-service/internal/model"
)

func (s *PvzService) AddNewPvz(ctx context.Context, city string) (*model.Pvz, error) {
	if err := s.validateCity(city); err != nil {
		return nil, err
	}

	idPvz, err := s.repo.CreatePvz(ctx, city)
	if err != nil {
		return nil, err
	}

	pvz, err := s.repo.GetPvzByID(ctx, idPvz)
	if err != nil {
		return nil, err
	}

	return pvz, nil
}

func (s *PvzService) validateCity(city string) error {
	switch city {
	case "Москва":
		return nil
	case "Санкт-Петербург":
		return nil
	case "Казань":
		return nil
	}
	return fmt.Errorf("invalid city: %s", city)
}
