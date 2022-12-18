package models

import "errors"

type Teacher struct {
	Id       int    `json:"id" gorm:"primaryKey"`
	Fullname string `json:"fullname"`
	UserId   int    `json:"-"`
	User     User   `json:"-" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Username string `json:"username" gorm:"->"`
}

var TeacherNotFound = errors.New("Teacher not found")
