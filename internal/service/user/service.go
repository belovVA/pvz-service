package user

import (
	"context"

	"github.com/google/uuid"
	"pvz-service/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (uuid.UUID, error)
	Get(ctx context.Context, email, password string) (*model.User, error)
}

type UserService struct {
	userRepository UserRepository
	jwtSecret      string
}

func NewUserService(
	repo UserRepository, jwt string,
) *UserService {
	return &UserService{
		userRepository: repo,
		jwtSecret:      jwt,
	}
}
