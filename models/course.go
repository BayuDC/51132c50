package models

import (
	"errors"

	"gorm.io/gorm"
)

type Course struct {
	Id        int      `json:"id" grom:"primaryKey"`
	Name      string   `json:"name"`
	TeacherId *int     `json:"-"`
	Teacher   *Teacher `json:"teacher" gorm:"foreignKey:TeacherId;constraint:OnDelete:SET NULL"`
}

var CourseTeacherNotFound = errors.New("Teacher not found")

func (c *Course) BeforeSave(tx *gorm.DB) (err error) {
	if c.TeacherId == nil {
		return
	}

	var teacher Teacher
	var result = tx.Take(&teacher, *c.TeacherId)
	if result.RowsAffected == 0 {
		err = CourseTeacherNotFound
	} else {
		err = result.Error
	}
	c.Teacher = &teacher

	return
}
