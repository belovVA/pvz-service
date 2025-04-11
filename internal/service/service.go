package service

import (
	"pvz-service/internal/service/auth"
	"pvz-service/internal/service/product"
	"pvz-service/internal/service/pvz"
	"pvz-service/internal/service/reception"
)

type Repository interface {
	auth.UserRepository
	pvz.PvzRepository
	reception.ReceptionRepository
	product.ProductRepository
}

type Service struct {
	*auth.AuthService
	*pvz.PvzService
	*reception.ReceptionService
	*product.ProductService
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		AuthService:      auth.NewAuthService(repo, jwtSecret),
		PvzService:       pvz.NewPvzService(repo),
		ReceptionService: reception.NewReceptionService(repo),
		ProductService:   product.NewProductService(repo, repo),
	}
}
