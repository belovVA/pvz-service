package model

import "time"

type PvzInfoQuery struct {
	StartDate time.Time
	EndDate   time.Time
	Page      int
	Limit     int
}

func (q *PvzInfoQuery) SetDefaults() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 30 {
		q.Limit = 10
	}
}
