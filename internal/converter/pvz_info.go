package converter

import (
	"time"

	"pvz-service/internal/handler/dto"
	"pvz-service/internal/model"
)

func ToPvzInfoQueryFromPvzInfoResponse(req *dto.PvzInfoRequest) (*model.PvzInfoQuery, error) {
	const layout = time.RFC3339

	start, err := time.Parse(layout, req.StartDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(layout, req.EndDate)
	if err != nil {
		return nil, err
	}

	ans := &model.PvzInfoQuery{
		StartDate: start,
		EndDate:   end,
		Page:      req.Page,
		Limit:     req.Limit,
	}
	ans.SetDefaults()
	return ans, nil

}
