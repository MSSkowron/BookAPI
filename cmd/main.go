package main

import (
	"log"

	"github.com/MSSkowron/BookRESTAPI/api"
	"github.com/MSSkowron/BookRESTAPI/storage"
)

func main() {
	storage, err := storage.NewPostgresSQLStorage()
	if err != nil {
		log.Fatalln(err)
	}

	api.NewBookRESTAPIServer(":8080", storage).Run()
}
