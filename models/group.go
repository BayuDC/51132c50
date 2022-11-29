package models

type Group struct {
	Id   int    `json:"id" grom:"primaryKey"`
	Name string `json:"name"`
}
