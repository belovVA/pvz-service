package auth

import (
	"context"

	"github.com/google/uuid"
	"pvz-service/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (uuid.UUID, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
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
