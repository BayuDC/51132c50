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

	count := int64(0)
	err = tx.Model(&Teacher{}).
		Where("id = ?", c.TeacherId).
		Count(&count).
		Error
	if count == 0 {
		err = CourseTeacherNotFound
	}

	return
}
