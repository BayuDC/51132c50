package models

type Student struct {
	Id       int    `json:"id" gorm:"primaryKey"`
	Fullname string `json:"fullname"`
	UserId   int    `json:"-"`
	User     User   `json:"-" gorm:"foreignKey:UserId;constraint:OnDelete:SET NULL"`
	Username string `json:"username" gorm:"->"`
}
