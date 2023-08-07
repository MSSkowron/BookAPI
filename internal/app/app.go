package app

import (
	"flag"
	"fmt"

	"github.com/MSSkowron/BookRESTAPI/internal/config"
	"github.com/MSSkowron/BookRESTAPI/internal/database"
	"github.com/MSSkowron/BookRESTAPI/internal/transport/rest"
)

func Run() error {
	configFileFlag := flag.String("configFile", "./configs/default_config.env", "path to a configuration file")
	flag.Parse()

	config, err := config.LoadConfig(*configFileFlag)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	database, err := database.NewPostgresqlDatabase(config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	if err := rest.NewServer(config.HTTPServerListenAddress, config.TokenSecret, config.TokenDuration, database).Run(); err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil
}
