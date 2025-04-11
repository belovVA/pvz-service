package modelRepo

import (
	"time"

	"github.com/google/uuid"
)

type Pvz struct {
	ID               uuid.UUID `db:"id"`
	RegistrationDate time.Time `db:"registration_date"`
	City             string    `db:"city"`
}
