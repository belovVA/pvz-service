package model

import (
	"time"

	"github.com/google/uuid"
)

type Pvz struct {
	ID               uuid.UUID
	RegistrationDate time.Time
	City             string
	Receptions       []Reception
}
