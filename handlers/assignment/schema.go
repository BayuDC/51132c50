package assignment

import "time"

type SetScheduleSchema struct {
	Group   int        `json:"group" binding:"required"`
	OpenAt  *time.Time `json:"open"`
	CloseAt *time.Time `json:"close"`
}
