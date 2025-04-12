package info

import (
	"pvz-service/internal/service/product"
	"pvz-service/internal/service/pvz"
	"pvz-service/internal/service/reception"
)

type InfoService struct {
	productRepository   product.ProductRepository
	receptionRepository reception.ReceptionRepository
	pvzRepository       pvz.PvzRepository
}

func NewInfoService(repoProduct product.ProductRepository, repoRepository reception.ReceptionRepository, repoPvz pvz.PvzRepository) *InfoService {
	return &InfoService{
		productRepository:   repoProduct,
		receptionRepository: repoRepository,
		pvzRepository:       repoPvz,
	}
}
