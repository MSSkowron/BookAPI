package main

import (
	"log"

	"github.com/MSSkowron/BookRESTAPI/api"
	"github.com/MSSkowron/BookRESTAPI/storage"
)

func main() {
	storage, err := storage.NewPostgresSQLStorage()
	if err != nil {
		log.Fatalf("error while creating storage: %s", err.Error())
	}

	api.NewGoBookAPIServer(":8080", storage).Run()
}
