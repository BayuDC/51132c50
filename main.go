package main

import (
	"tink/models"
	"tink/server"
)

func main() {
	db := models.ConnectDb()
	s := server.New(db)
	s.Run()
}
