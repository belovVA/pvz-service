package auth

import (
	"context"
	"fmt"
	"time"

	"pvz-service/internal/model"
	"pvz-service/pkg/hash_password"
	"pvz-service/pkg/jwtutils"
)

func (s *AuthService) Authenticate(ctx context.Context, user model.User) (string, error) {

	hashPass, err := hash_password.HashPassword(user.Password)
	if err != nil {
		return "", err
	}

	current, err := s.userRepository.GetByEmailAndPass(ctx, user.Email, hashPass)
	if err != nil {
		return "", err
	}

	token, err := s.generateJWT(current.ID.String(), current.Role)

	return token, nil
}

func (s *AuthService) generateJWT(userID string, role string) (string, error) {
	claims := map[string]interface{}{
		"user_id": userID,
		"role":    role,
	}

	token, err := jwtutils.Generate(claims, 24*time.Hour, s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token")
	}

	return token, nil
}

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

	return nil, err
}
