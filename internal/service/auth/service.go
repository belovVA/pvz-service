package auth

import (
	"context"

	"github.com/google/uuid"
	"pvz-service/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (uuid.UUID, error)
	GetByEmailAndPass(ctx context.Context, email, password string) (*model.User, error)
}

type AuthService struct {
	userRepository UserRepository
	jwtSecret      string
}

func NewAuthService(
	repo UserRepository, jwt string,
) *AuthService {
	return &AuthService{
		userRepository: repo,
		jwtSecret:      jwt,
	}
}
