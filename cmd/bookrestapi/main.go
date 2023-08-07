package main

import (
	"flag"
	"log"

	"github.com/MSSkowron/BookRESTAPI/internal/config"
	"github.com/MSSkowron/BookRESTAPI/internal/server"
	"github.com/MSSkowron/BookRESTAPI/internal/storage"
)

func main() {
	configFileFlag := flag.String("configFile", "./configs/default_config.env", "path to a configuration file")
	flag.Parse()

	config, err := config.LoadConfig(*configFileFlag)
	if err != nil {
		log.Fatalf("Error while loading config: %s", err)
	}

	storage, err := storage.NewPostgresStorage(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Error while creating storage: %s", err)
	}

	if err := server.NewServer(config.HTTPServerListenAddress, config.TokenSecret, config.TokenDuration, storage).Run(); err != nil {
		log.Fatalf("Error while running server: %s", err)
	}
}
