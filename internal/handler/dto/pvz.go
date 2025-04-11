package dto

import "time"

type CreatePvzRequest struct {
	City string `json:"city" validate:"required"`
}

type CreatePvzResponse struct {
	ID               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}
