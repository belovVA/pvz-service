package dto

import "time"

type CreateReceptionRequest struct {
	PvzID string `json:"pvzId" validate:"required"`
}

type CreateReceptionResponse struct {
	ID       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PvzID    string    `json:"pvzId"`
	Status   string    `json:"status"`
}
