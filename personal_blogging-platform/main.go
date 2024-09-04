package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	store, err := newPostgresConnection()
	if err != nil {
		log.Fatal("Can not connect to database", err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	server := newAPIServer(":8080", store)
	server.Run()
}
