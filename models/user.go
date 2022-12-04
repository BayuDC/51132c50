package models

import (
	"errors"

	"gorm.io/gorm"
)

type role string

const (
	RoleAdmin   role = "admin"
	RoleTeacher role = "teacher"
	RoleStudent role = "student"
)

type User struct {
	Id       int    `json:"id" gorm:"primaryKey"`
	Fullname string `json:"fullname" gorm:"->"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"-"`
	Role     role   `json:"role" gorm:"type:role"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	count := int64(0)
	err = tx.Model(&User{}).
		Where("username = ?", u.Username).
		Count(&count).
		Error
	if count > 0 {
		err = errors.New("Username is taken")
	}
	return
}
