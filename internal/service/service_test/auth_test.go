package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"pvz-service/internal/model"
	"pvz-service/internal/service"
	"pvz-service/internal/service/mocks"
)

func TestAuthService_Registration(t *testing.T) {
	ctx := context.Background()

	password := "securepassword"
	//hashedPassword, _ := pkg.HashPassword(password) // желательно использовать стабильную фиктивную реализацию

	testUser := model.User{
		Email:    "test@example.com",
		Password: password,
		Role:     "user",
	}

	tests := []struct {
		name              string
		setupMocks        func(repo *mocks.UserRepository)
		expectedErr       error
		expectedUserEmail string
	}{
		{
			name: "successful registration",
			setupMocks: func(repo *mocks.UserRepository) {
				repo.On("GetUserByEmail", mock.Anything, testUser.Email).
					Return(nil, errors.New("not found"))

				repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
					return u.Email == testUser.Email && u.Password != password
				})).Return(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), nil)
			},
			expectedErr:       nil,
			expectedUserEmail: testUser.Email,
		},
		{
			name: "user already exists",
			setupMocks: func(repo *mocks.UserRepository) {
				repo.On("GetUserByEmail", mock.Anything, testUser.Email).
					Return(&model.User{}, nil) // пользователь найден — ошибка
			},
			expectedErr: errors.New("user already exist"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)
			tt.setupMocks(mockRepo)

			authService := service.NewAuthService(mockRepo, "test")

			result, err := authService.Registration(ctx, testUser)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedUserEmail, result.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
func TestAuthService_DummyAuth(t *testing.T) {
	hashedPassModerator, _ := bcrypt.GenerateFromPassword([]byte(service.ModeratorRole), bcrypt.DefaultCost)
	hashedPassEmployee, _ := bcrypt.GenerateFromPassword([]byte(service.EmployeeRole), bcrypt.DefaultCost)

	moderatorID := uuid.New()

	tests := []struct {
		name       string
		role       string
		email      string
		setupMocks func(repo *mocks.UserRepository)
		wantErr    bool
	}{
		{
			name:  "moderator login success",
			role:  service.ModeratorRole,
			email: service.ModeratorEmail,
			setupMocks: func(repo *mocks.UserRepository) {
				repo.On("GetUserByEmail", mock.Anything, service.ModeratorEmail).
					Return(&model.User{
						ID:       moderatorID,
						Email:    service.ModeratorEmail,
						Password: string(hashedPassModerator),
						Role:     service.ModeratorRole,
					}, nil).Once()
			},
			wantErr: false,
		},
		{
			name:  "employee login success",
			role:  service.EmployeeRole,
			email: service.EmployeeEmail,
			setupMocks: func(repo *mocks.UserRepository) {
				repo.On("GetUserByEmail", mock.Anything, service.EmployeeEmail).
					Return(&model.User{
						ID:       moderatorID,
						Email:    service.EmployeeEmail,
						Password: string(hashedPassEmployee),
						Role:     service.EmployeeEmail,
					}, nil).Once()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)
			tt.setupMocks(mockRepo)

			authService := service.NewAuthService(mockRepo, "test")

			userDummy := model.User{
				Email:    tt.email,
				Password: tt.role,
				Role:     tt.role,
			}

			token, err := authService.DummyAuth(context.Background(), userDummy)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, token)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
