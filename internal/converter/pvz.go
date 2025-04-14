package converter

import (
	"github.com/google/uuid"
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/model"
)

func ToCreatePvzResponseFromPvz(pvz *model.Pvz) *dto.PvzResponse {
	return &dto.PvzResponse{
		ID:               pvz.ID.String(),
		RegistrationDate: pvz.RegistrationDate,
		City:             pvz.City,
	}
}

func ToPvzFromCreatePvzRequest(dto *dto.CreatePvzRequest) *model.Pvz {
	return &model.Pvz{
		ID:   uuid.Nil,
		City: dto.City,
	}
}

func ToPvzFromIDRequest(IDString string) (*model.Pvz, error) {
	id, err := uuid.Parse(IDString)
	return &model.Pvz{ID: id}, err
}
