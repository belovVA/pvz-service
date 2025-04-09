package converter

import (
	"pvz-service/internal/model"
	modelRepo "pvz-service/internal/repository/pgdb/model"
)

func ToUserFromUserRepo(user *modelRepo.User) *model.User {
	return &model.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}
}
