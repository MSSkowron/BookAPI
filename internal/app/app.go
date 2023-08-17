package app

import (
	"flag"
	"fmt"

	"github.com/MSSkowron/BookRESTAPI/internal/api"
	"github.com/MSSkowron/BookRESTAPI/internal/config"
	"github.com/MSSkowron/BookRESTAPI/internal/database"
	"github.com/MSSkowron/BookRESTAPI/internal/services"
)

// Run runs the BookRESTAPI application.
// It loads configuration, creates database connection, creates services and runs the server.
// It returns an error if any of the steps fails.
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

	tokenService := services.NewTokenService(config.TokenSecret, config.TokenDuration)
	userService := services.NewUserService(database, tokenService)
	bookService := services.NewBookService(database)

	if err := api.NewServer(userService, bookService, tokenService, api.WithAddress(config.HTTPServerListenAddress)).ListenAndServe(); err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil
}
