package converter

import (
	"github.com/google/uuid"
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/model"
)

func ToUserFromCreateUserRequest(user *dto.CreateUserRequest) *model.User {
	return &model.User{
		ID:       uuid.Nil,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}
}

func ToCreateUserResponseFromUser(user *model.User) *dto.CreateUserResponse {
	return &dto.CreateUserResponse{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role,
	}
}

func ToUserFromLoginUserRequest(user *dto.LoginUserRequest) *model.User {
	return &model.User{
		ID:       uuid.Nil,
		Email:    user.Email,
		Password: user.Password,
		Role:     "",
	}
}
