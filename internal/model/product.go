package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	DateTime    time.Time
	TypeProduct string
	ReceptionID uuid.UUID
}
