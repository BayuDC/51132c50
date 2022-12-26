package models

import "errors"

type Group struct {
	Id   int    `json:"id" grom:"primaryKey"`
	Name string `json:"name"`
}

var GroupNotFound = errors.New("Group not found")
