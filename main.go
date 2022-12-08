package main

import "github.com/MSSkowron/GoBankAPI/api"

func main() {
	api.NewGoBookAPIServer(":8080").Run()
}
