package service

import (
	"pvz-service/internal/service/auth"
)

type Repository interface {
	auth.UserRepository
}

type Service struct {
	*auth.AuthService
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		AuthService: auth.NewAuthService(repo, jwtSecret),
	}
}
