package main

import (
	"log"

	"github.com/MSSkowron/GoBankAPI/api"
	"github.com/MSSkowron/GoBankAPI/storage"
)

func main() {
	storage, err := storage.NewPostgresSQLStorage()
	if err != nil {
		log.Fatalln("error while creating storage: " + err.Error())
	}

	api.NewGoBookAPIServer(":8080", storage).Run()
}
