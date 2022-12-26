package models

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	Id           int       `json:"id" grom:"primaryKey"`
	AssignmentId int       `json:"-"`
	GroupId      int       `json:"-"`
	OpenAt       time.Time `json:"open_at"`
	CloseAt      time.Time `json:"close_at"`
}

type Tabler interface {
	TableName() string
}

func (Schedule) TableName() string {
	return "assignment_schedules"
}

func (s *Schedule) BeforeSave(tx *gorm.DB) (err error) {
	var group Group
	var result = tx.Take(&group, *&s.GroupId)
	if result.RowsAffected == 0 {
		err = GroupNotFound
	} else {
		err = result.Error
	}

	return
}
