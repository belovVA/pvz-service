package user

import (
	"context"
	"fmt"

	"pvz-service/internal/model"
	"pvz-service/pkg/hash_password"
)

func (s *UserService) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	if err := ValidateRole(user.Role); err != nil {
		return nil, err
	}

	hashPass, err := hash_password.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	userID, err := s.userRepository.Create(ctx, &user)
	if err != nil {
		return nil, err
	}

	user.ID = userID
	user.Password = hashPass

	return &model.User{
		ID:       userID,
		Email:    user.Email,
		Password: hashPass,
		Role:     user.Role,
	}, nil
}

func ValidateRole(role string) error {
	switch role {
	case "employee":
		return nil
	case "moderator":
		return nil
	}
	return fmt.Errorf("invalid role")
}
