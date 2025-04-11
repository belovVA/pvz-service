package converter

import (
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/model"
)

func ToReceptionResponseFromReception(r *model.Reception) *dto.CreateReceptionResponse {
	return &dto.CreateReceptionResponse{
		ID:       r.ID.String(),
		DateTime: r.DateTime,
		PvzID:    r.PvzID.String(),
		Status:   r.Status(),
	}
}
