package models

import "errors"

type Student struct {
	Id       int    `json:"id" gorm:"primaryKey"`
	Fullname string `json:"fullname"`
	UserId   int    `json:"-"`
	User     User   `json:"-" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Username string `json:"username" gorm:"->"`
	GroupId  int    `json:"-"`
}

var StudentNotFound = errors.New("Student not found")
