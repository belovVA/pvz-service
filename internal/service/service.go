package service

import (
	"pvz-service/internal/service/user"
)

type Repository interface {
	user.UserRepository
}

type Service struct {
	*user.UserService
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		UserService: user.NewUserService(repo, jwtSecret),
	}
}
