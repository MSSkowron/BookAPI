package main

import (
	"log"

	"github.com/MSSkowron/BookRESTAPI/api"
	"github.com/MSSkowron/BookRESTAPI/config"
	"github.com/MSSkowron/BookRESTAPI/storage"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error while loading config: %s", err.Error())
	}

	storage, err := storage.NewPostgresSQLStorage(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("Error while creating storage: %s", err.Error())
	}

	api.NewBookRESTAPIServer(config.HTTPServerListenAddress, config.TokenSecret, config.TokenDuration, storage).Run()
}
