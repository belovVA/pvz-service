package service

import (
	"pvz-service/internal/service/auth"
	"pvz-service/internal/service/pvz"
)

type Repository interface {
	auth.UserRepository
	pvz.PvzRepository
}

type Service struct {
	*auth.AuthService
	*pvz.PvzService
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		AuthService: auth.NewAuthService(repo, jwtSecret),
		PvzService:  pvz.NewPvzService(repo),
	}
}
