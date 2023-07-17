package main

import (
	"flag"
	"log"

	"github.com/MSSkowron/BookRESTAPI/api"
	"github.com/MSSkowron/BookRESTAPI/config"
	"github.com/MSSkowron/BookRESTAPI/storage"
)

func main() {
	configFileFlag := flag.String("configFile", "./config.env", "path to a configuration file")
	flag.Parse()

	config, err := config.LoadConfig(*configFileFlag)
	if err != nil {
		log.Fatalf("Error while loading config: %s", err.Error())
	}

	storage, err := storage.NewPostgresSQLStorage(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Error while creating storage: %s", err.Error())
	}

	api.NewBookRESTAPIServer(config.HTTPServerListenAddress, config.TokenSecret, config.TokenDuration, storage).Run()
}
