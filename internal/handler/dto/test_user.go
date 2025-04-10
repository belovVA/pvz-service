package dto

type TestUserRequest struct {
	Role string `json:"role" binding:"required"`
}
