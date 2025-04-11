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

	current, err := s.userRepository.GetUserByEmail(ctx, user.Email)
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
