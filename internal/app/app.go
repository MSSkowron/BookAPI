package app

import (
	"flag"
	"fmt"

	"github.com/MSSkowron/BookRESTAPI/internal/api"
	"github.com/MSSkowron/BookRESTAPI/internal/config"
	"github.com/MSSkowron/BookRESTAPI/internal/database"
	"github.com/MSSkowron/BookRESTAPI/internal/services"
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
		return fmt.Errorf("failed to create database: %w", err)
	}

	userService := services.NewUserService(database, config.TokenSecret, config.TokenDuration)
	bookService := services.NewBookService(database)

	if err := api.NewServer(config.HTTPServerListenAddress, userService, bookService).Run(); err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil
}
