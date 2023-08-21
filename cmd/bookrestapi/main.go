package main

import (
	"github.com/MSSkowron/BookRESTAPI/internal/app"
	"github.com/MSSkowron/BookRESTAPI/pkg/logger"
)

func main() {
	if err := app.Run(); err != nil {
		logger.Errorln(err)
	}
}
