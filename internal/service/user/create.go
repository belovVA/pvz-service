package user

import (
	"context"

	"pvz-service/internal/model"
)

func (s *userServ) Register(ctx context.Context, user model.User) (*model.User, error) {
	// TODO Hash Password
	// TODO CheckRole and move this func to user/service prob
	//user.Password =

	userID, err := s.userRepository.Create(ctx, &user)
	if err != nil {
		return nil, err
	}
	user.ID = userID
	return &user, nil
}
