package converter

import (
	"pvz-service/internal/model"
	modelRepo "pvz-service/internal/repository/pgdb/model"
)

func ToPvzFromPvzRepo(pvz *modelRepo.Pvz) *model.Pvz {
	return &model.Pvz{
		ID:               pvz.ID,
		RegistrationDate: pvz.RegistrationDate,
		City:             pvz.City,
	}
}
