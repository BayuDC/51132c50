package models

import (
	"gorm.io/gorm"
)

type Course struct {
	Id        int       `json:"id" grom:"primaryKey"`
	Name      string    `json:"name"`
	TeacherId *int      `json:"-"`
	Teacher   *Teacher  `json:"teacher" gorm:"foreignKey:TeacherId;constraint:OnDelete:SET NULL"`
	Students  []Student `json:"-" gorm:"many2many:student_courses;"`
}

func (c *Course) BeforeSave(tx *gorm.DB) (err error) {
	if c.TeacherId == nil {
		return
	}

	var teacher Teacher
	var result = tx.Take(&teacher, *c.TeacherId)
	if result.RowsAffected == 0 {
		err = TeacherNotFound
	} else {
		err = result.Error
	}
	c.Teacher = &teacher

	return
}
