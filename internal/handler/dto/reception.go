package dto

import "time"

type ReceptionRequest struct {
	PvzID string `json:"pvzId" validate:"required"`
}

type ReceptionResponse struct {
	ID       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PvzID    string    `json:"pvzId"`
	Status   string    `json:"status"`
}
