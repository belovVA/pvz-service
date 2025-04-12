package service

import (
	"context"

	"github.com/google/uuid"
	"pvz-service/internal/model"
)

type InfoService struct {
	productRepository   ProductRepository
	receptionRepository ReceptionRepository
	pvzRepository       PvzRepository
}

func NewInfoService(repoProduct ProductRepository, repoRepository ReceptionRepository, repoPvz PvzRepository) *InfoService {
	return &InfoService{
		productRepository:   repoProduct,
		receptionRepository: repoRepository,
		pvzRepository:       repoPvz,
	}
}

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

	recepMap := s.getMapReceptionsOrderByPvz(receptions)

	if query.StartDate.IsZero() && query.EndDate.IsZero() {
		idList, err := s.pvzRepository.GetIDListPvz(ctx)
		if err != nil {
			return nil, err
		}
		for _, id := range idList {
			if _, ok := recepMap[id]; !ok {
				recepMap[id] = []model.Reception{}
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

	return s.rangeAnswerByPagination(res, query.Page, query.Limit), nil
}

func (s *InfoService) getMapReceptionsOrderByPvz(receptions []model.Reception) map[uuid.UUID][]model.Reception {
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
	return recepMap
}

func (s *InfoService) rangeAnswerByPagination(res []*model.Pvz, page int, limit int) []*model.Pvz {
	start := (page - 1) * limit
	if start >= len(res) {
		return []*model.Pvz{}
	}

	end := start + limit
	if end > len(res) {
		end = len(res)
	}

	return res[start:end]
}
