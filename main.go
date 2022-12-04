package main

import (
	"tink/models"
	"tink/server"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db := models.ConnectDb()
	s := server.New(db)
	s.Run()
}
