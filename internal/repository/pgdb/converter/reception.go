package converter

import (
	"pvz-service/internal/model"
	modelRepo "pvz-service/internal/repository/pgdb/model"
)

func ToReceptionFromReceptionRepo(reception *modelRepo.Reception) *model.Reception {
	return &model.Reception{
		ID:       reception.ID,
		DateTime: reception.DateTime,
		Products: nil,
		IsClosed: reception.IsClosedStatus,
		PvzID:    reception.PvzID,
	}
}
