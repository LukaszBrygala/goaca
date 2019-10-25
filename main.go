package main

import (
	"log"

	"academy/db"
	"academy/gateways"
	"academy/server"
)

func main() {
	db, err := db.Connect(db.Config{
		Host: "localhost",
		Port: 5432,
		User: "local",
		Pass: "pass",
		Name: "dc",
	})
	if err != nil {
		log.Fatal(err)
	}

	repo := gateways.NewRepository(db)

	s := server.New(repo)
	s.Run(8080)
}
