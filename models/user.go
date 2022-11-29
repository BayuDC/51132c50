package models

type role string

const (
	RoleAdmin   role = "admin"
	RoleTeacher role = "teacher"
	RoleStudent role = "student"
)

type User struct {
	Id       int    `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Role     role   `json:"role" gorm:"type:role"`
}
