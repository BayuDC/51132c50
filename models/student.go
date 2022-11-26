package models

type Student struct {
	Id   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}
