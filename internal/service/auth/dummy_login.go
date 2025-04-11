package auth

import (
	"context"
	"fmt"

	"pvz-service/internal/model"
)

func (s *AuthService) DummyAuth(ctx context.Context, role string) (string, error) {
	user, err := getTestUserByRole(role)
	if err != nil {
		return "", err
	}

	token, err := s.Authenticate(ctx, *user)

	// Если данный тестовый пользователь не был найден
	// Сработает только при первом вызове, либо при краше БД
	// Попробуем его зарегать, а потом заново зайти
	// Если получаем ошибку, то либо сменили пароль учетки, либо краш БД
	if err != nil {
		_, err = s.Registration(ctx, *user)
		if err != nil {
			return "", err
		}

		token, err = s.Authenticate(ctx, *user)
		if err != nil {
			return "", err
		}
	}
	return token, nil

}

func getTestUserByRole(role string) (*model.User, error) {
	switch role {
	case "employee":
		return &model.User{
			Email:    "employee@test.com",
			Password: role,
			Role:     "employee",
		}, nil

	case "moderator":
		return &model.User{
			Email:    "moderator@test.com",
			Password: role,
			Role:     "moderator",
		}, nil
	}

	return nil, fmt.Errorf("forbidden role")
}
