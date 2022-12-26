package models

type Assignment struct {
	Id          int    `json:"id" grom:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CourseId    int    `json:"-"`
	Course      Course `json:"-"`
	Type        string `json:"type"`
}
