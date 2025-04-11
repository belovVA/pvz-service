package dto

import "time"

type CreateProductRequest struct {
	TypeProduct string `json:"type" validate:"required"`
	PvzID       string `json:"pvzId" validate:"required"`
}

type CreateProductResponse struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	TypeProduct string    `json:"type"`
	ReceptionID string    `json:"receptionId"`
}
