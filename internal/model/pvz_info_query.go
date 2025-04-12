package model

import "time"

type PvzInfoQuery struct {
	StartDate time.Time
	EndDate   time.Time
	Page      int
	Limit     int
}
