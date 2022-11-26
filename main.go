package main

import (
	"tink/server"
)

func main() {
	s := server.New()
	s.Run()
}
