package main

import (
	"github.com/MSSkowron/BookRESTAPI/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
