package model

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	ID       uuid.UUID
	DateTime time.Time
	Products []Product
	IsClosed bool // true = "closed", false = "in_progress"
	PvzID    uuid.UUID
}

func (r *Reception) Status() string {
	if r.IsClosed {
		return "closed"
	}
	return "in_progress"
}
