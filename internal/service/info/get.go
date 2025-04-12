package info

import (
	"context"

	"github.com/google/uuid"
	"pvz-service/internal/model"
)

func (s *InfoService) GetInfoPvz(ctx context.Context, query *model.PvzInfoQuery) ([]*model.Pvz, error) {
	receptions, err := s.receptionRepository.GetReceptionsSliceWithTimeRange(ctx, query.StartDate, query.EndDate)
	if err != nil {
		return nil, err
	}

	for i, _ := range receptions {
		receptions[i].Products, err = s.productRepository.GetProductSliceByReceptionID(ctx, receptions[i].ID)
		if err != nil {
			return nil, err
		}
	}

	// Распределяем receptions согласно их PVZ

	recepMap := make(map[uuid.UUID][]model.Reception, len(receptions))
	for _, recep := range receptions {
		key := recep.PvzID
		if _, ok := recepMap[key]; ok {
			recepMap[key] = append(recepMap[key], recep)
		} else {

			recepMap[key] = make([]model.Reception, 0, 10)
			recepMap[key] = append(recepMap[key], recep)

		}
	}

	if query.StartDate.IsZero() && query.EndDate.IsZero() {
		idList, err := s.pvzRepository.GetIDListPvz(ctx)
		if err != nil {
			return nil, err
		}
		for _, id := range idList {
			if _, ok := recepMap[id]; !ok {
				recepMap[id] = make([]model.Reception, 0, 10)
			}
		}
	}

	res := make([]*model.Pvz, 0, len(recepMap))

	// Получаем данные для каждого ПВЗ и заносим в res
	for k, _ := range recepMap {
		pvz, err := s.pvzRepository.GetPvzByID(ctx, k)
		if err != nil {
			return nil, err
		}
		pvz.Receptions = recepMap[k]
		res = append(res, pvz)
	}

	start := (query.Page - 1) * query.Limit
	if start >= len(res) {
		return []*model.Pvz{}, nil
	}
	end := start + query.Limit
	if end > len(res) {
		end = len(res)
	}

	return res[start:end], nil
}
