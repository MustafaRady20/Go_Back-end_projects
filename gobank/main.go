package main

import (
	"fmt"
	"log"
)

func main() {
	store, err := newPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", store)
	server := newAPIServer(":8080", store)
	server.Run()
}
