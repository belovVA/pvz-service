package converter

import (
	"github.com/google/uuid"
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/model"
)

func ToReceptionResponseFromReception(r *model.Reception) *dto.ReceptionResponse {
	return &dto.ReceptionResponse{
		ID:       r.ID.String(),
		DateTime: r.DateTime,
		PvzID:    r.PvzID.String(),
		Status:   r.Status(),
	}
}

func ToReceptionFromReceptionRequest(dto *dto.ReceptionRequest) (*model.Reception, error) {
	pvzID, err := uuid.Parse(dto.PvzID)
	return &model.Reception{
		ID:    uuid.Nil,
		PvzID: pvzID,
	}, err
}

func ToReceptionFromPvzIDRequest(IDString string) (*model.Reception, error) {
	id, err := uuid.Parse(IDString)
	return &model.Reception{
		PvzID: id,
	}, err
}
