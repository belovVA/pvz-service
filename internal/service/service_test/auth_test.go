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
					return u.Email == testUser.Email && u.Password != "" // пароль будет захеширован
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

func TestAuthService_Authenticate(t *testing.T) {
	type fields struct {
		userRepo    *mocks.UserRepository
		generateJWT func(userID, role string) (string, error)
	}

	type args struct {
		user model.User
	}

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	userID := uuid.New()

	tests := []struct {
		name      string
		args      args
		mockSetup func(repo *mocks.UserRepository)
		fields    fields
		wantToken string
		wantErr   string
	}{
		{
			name: "successful auth",
			args: args{
				user: model.User{
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(&model.User{
					ID:       userID,
					Email:    "test@example.com",
					Password: string(hashedPass),
					Role:     "user",
				}, nil)
			},
			fields: fields{
				userRepo: &mocks.UserRepository{},
				generateJWT: func(userID, role string) (string, error) {
					return "mocked-jwt-token", nil
				},
			},
			wantToken: "mocked-jwt-token",
			wantErr:   "",
		},
		{
			name: "user not found",
			args: args{
				user: model.User{Email: "notfound@example.com", Password: "somepass"},
			},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetUserByEmail", mock.Anything, "notfound@example.com").Return(nil, errors.New("not found"))
			},
			fields: fields{
				userRepo: &mocks.UserRepository{},
				generateJWT: func(userID, role string) (string, error) {
					return "", nil
				},
			},
			wantToken: "",
			wantErr:   "user not found",
		},
		{
			name: "invalid password",
			args: args{
				user: model.User{Email: "test@example.com", Password: "wrongpass"},
			},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(&model.User{
					ID:       userID,
					Email:    "test@example.com",
					Password: string(hashedPass),
					Role:     "user",
				}, nil)
			},
			fields: fields{
				userRepo: &mocks.UserRepository{},
				generateJWT: func(userID, role string) (string, error) {
					return "", nil
				},
			},
			wantToken: "",
			wantErr:   "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := new(mocks.UserRepository)
			tt.mockSetup(mockRepo)

			authService := service.NewAuthService(mockRepo, "test")

			_, err := authService.Authenticate(context.Background(), tt.args.user)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
