package models

import (
	"tink/middlewares"

	"gorm.io/gorm"
)

type Course struct {
	Id          int          `json:"id" grom:"primaryKey"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	TeacherId   *int         `json:"-"`
	Teacher     *Teacher     `json:"teacher" gorm:"foreignKey:TeacherId;constraint:OnDelete:SET NULL"`
	Students    []Student    `json:"-" gorm:"many2many:student_courses;"`
	Assignments []Assignment `json:"-"`
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

func (c *Course) Check(tx *gorm.DB, user *middlewares.User) bool {
	switch user.Role {
	case "teacher":
		if c.TeacherId == nil {
			return false
		} else if *c.TeacherId != user.Userable {
			return false
		}
	case "student":
		if tx.Model(&c).
			Where("students.id = ?", user.Userable).
			Association("Students").
			Count() == 0 {
			return false
		}
	default:
		return false
	}

	return true
}
