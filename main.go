package main

import (
	"academy/server"
)

func main() {
	s := server.New()
	s.Run(8080)
}
