package modelRepo

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `db:"id"`
	DateTime    time.Time `db:"date_time"`
	TypeProduct string    `db:"type_product"`
	ReceptionID uuid.UUID `db:"reception_id, foreign key"`
}
