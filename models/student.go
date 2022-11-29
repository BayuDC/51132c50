package models

type Student struct {
	Id     int    `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	UserId int    `json:"user_id"`
	User   User   `json:"user" gorm:"references:UserId"`
}
