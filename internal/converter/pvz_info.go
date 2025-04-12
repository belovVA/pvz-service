package converter

import (
	"time"

	"pvz-service/internal/handler/dto"
	"pvz-service/internal/model"
)

func ToPvzInfoQueryFromPvzInfoResponse(req *dto.PvzInfoRequest) (*model.PvzInfoQuery, error) {
	const layout = time.RFC3339
	start, err := time.Parse(layout, req.StartDate)
	if err != nil && req.StartDate != "" {
		return nil, err
	}

	end, err := time.Parse(layout, req.EndDate)
	if err != nil && req.EndDate != "" {
		return nil, err
	}

	ans := &model.PvzInfoQuery{
		StartDate: start,
		EndDate:   end,
		Page:      req.Page,
		Limit:     req.Limit,
	}
	ans.SetDefaults()
	//log.Println(ans, ans.StartDate.IsZero())

	return ans, nil

}

func ToPvzInfoResponseList(pvzList []*model.Pvz) []dto.SinglePvzInfoResponse {
	responseList := make([]dto.SinglePvzInfoResponse, 0, len(pvzList))

	for _, pvz := range pvzList {
		pvzResp := ToCreatePvzResponseFromPvz(pvz)

		receptions := make([]dto.ReceptionInfo, 0, len(pvz.Receptions))
		for _, rec := range pvz.Receptions {
			recResp := ToReceptionResponseFromReception(&rec)

			products := make([]dto.ProductResponse, 0, len(rec.Products))
			for _, prod := range rec.Products {
				prodResp := ToProductResponseFromProduct(&prod)
				products = append(products, *prodResp)
			}

			receptions = append(receptions, dto.ReceptionInfo{
				ReceptionData: *recResp,
				Products:      products,
			})
		}

		responseList = append(responseList, dto.SinglePvzInfoResponse{
			PvzData:    *pvzResp,
			Receptions: receptions,
		})
	}

	return responseList
}
