package auth

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"pvz-service/internal/model"
	"pvz-service/pkg/jwtutils"
)

func (s *AuthService) Authenticate(ctx context.Context, user model.User) (string, error) {

	current, err := s.userRepository.GetByEmail(ctx, user.Email)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(current.Password), []byte(user.Password)); err != nil {
		return "", fmt.Errorf("invalid email or password")
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

	return nil, fmt.Errorf("forbidden role")
}
