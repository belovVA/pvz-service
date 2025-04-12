package converter

import (
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
