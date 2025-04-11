package auth

import (
	"context"
	"fmt"

	"pvz-service/internal/model"
	"pvz-service/internal/service/pkg"
)

func (s *AuthService) Registration(ctx context.Context, user model.User) (*model.User, error) {
	if err := s.validateRole(user.Role); err != nil {
		return nil, err
	}

	hashPass, err := pkg.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashPass

	if _, err := s.userRepository.GetUserByEmail(ctx, user.Email); err == nil {
		return nil, fmt.Errorf("user already exist")
	}

	userID, err := s.userRepository.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:       userID,
		Email:    user.Email,
		Password: hashPass,
		Role:     user.Role,
	}, nil
}

func (s *AuthService) validateRole(role string) error {
	switch role {
	case "employee":
		return nil
	case "moderator":
		return nil
	}
	return fmt.Errorf("invalid role: %s", role)
}
