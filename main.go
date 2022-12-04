package main

import (
	"log"
	"tink/models"
	"tink/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := models.ConnectDb()
	s := server.New(db)
	s.Run()
}
