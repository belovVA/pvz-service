package converter

import (
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/model"
)

func ToCreatePvzResponseFromPvz(pvz *model.Pvz) *dto.CreatePvzResponse {
	return &dto.CreatePvzResponse{
		ID:               pvz.ID.String(),
		RegistrationDate: pvz.RegistrationDate,
		City:             pvz.City,
	}
}
