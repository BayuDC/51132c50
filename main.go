package main

import (
	"tink/core"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db := core.ConnectDb()
	app := core.CreateApp(db)
	app.Run()
}
