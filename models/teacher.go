package models

type Teacher struct {
	Id     int    `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	UserId int    `json:"user_id"`
}
