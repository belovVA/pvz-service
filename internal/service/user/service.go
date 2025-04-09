package user

import "pvz-service/internal/repository"

type userServ struct {
	userRepository repository.UserRepository
}

func NewUserservice(
	userRepository repository.UserRepository,
) *userServ {
	return &userServ{
		userRepository: userRepository,
	}
}
