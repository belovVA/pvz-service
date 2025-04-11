package modelRepo

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	ID             uuid.UUID `db:"id"`
	DateTime       time.Time `db:"date_time"`
	IsClosedStatus bool      `db:"is_closed"`
	PvzID          uuid.UUID `db:"pvz_id, foreign key"`
}
