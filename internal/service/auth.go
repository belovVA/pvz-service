package service

import (
	"context"
	"fmt"
	"time"

	"pvz-service/internal/model"
	"pvz-service/internal/service/pkg/hash"
	"pvz-service/pkg/jwtutils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
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

func (s *AuthService) Registration(ctx context.Context, user model.User) (*model.User, error) {
	hashPass, err := hash.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash pass")
	}

	user.Password = hashPass

	if _, err := s.userRepository.GetUserByEmail(ctx, user.Email); err == nil {
		return nil, fmt.Errorf("user already exist")
	}

	userID, err := s.userRepository.CreateUser(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user")
	}

	return &model.User{
		ID:       userID,
		Email:    user.Email,
		Password: hashPass,
		Role:     user.Role,
	}, nil
}

func (s *AuthService) Authenticate(ctx context.Context, user model.User) (string, error) {
	current, err := s.userRepository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(current.Password), []byte(user.Password)); err != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	token, err := s.generateJWT(current.ID.String(), current.Role)

	return token, err
}

func (s *AuthService) DummyAuth(ctx context.Context, user model.User) (string, error) {
	userDummy, err := getTestUserByRole(user.Role)
	if err != nil {
		return "", err
	}

	token, err := s.Authenticate(ctx, *userDummy)

	// Сработает только при первом вызове, либо при краше БД
	// Если данный тестовый пользователь не был найден
	// Попробуем его зарегать, а потом заново зайти
	// Если получаем ошибку, то либо сменили пароль учетки, либо краш БД
	if err != nil {
		_, err = s.Registration(ctx, *userDummy)
		if err != nil {
			return "", fmt.Errorf("failed to create test user")
		}

		token, err = s.Authenticate(ctx, *userDummy)
		if err != nil {
			return "", fmt.Errorf("failed to authenticate test user")
		}
	}

	return token, nil
}

func getTestUserByRole(role string) (*model.User, error) {
	const (
		employeeRole  = "employee"
		employeeEmail = "employee@test.com"

		moderatorRole  = "moderator"
		moderatorEmail = "moderator@test.com"
	)
	switch role {
	case employeeRole:
		return &model.User{
			Email:    employeeEmail,
			Password: role,
			Role:     employeeRole,
		}, nil

	case moderatorRole:
		return &model.User{
			Email:    moderatorEmail,
			Password: role,
			Role:     moderatorRole,
		}, nil
	}

	return nil, fmt.Errorf("forbidden role")
}

func (s *AuthService) generateJWT(userID string, role string) (string, error) {
	claims := map[string]interface{}{
		"userId": userID,
		"role":   role,
	}

	token, err := jwtutils.Generate(claims, 24*time.Hour, s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token")
	}

	return token, nil
}
