package dto

type TestUserRequest struct {
	Role string `json:"role" validate:"required"`
}
