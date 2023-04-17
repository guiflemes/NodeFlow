package main

import (
	"flowChart/server"
	"flowChart/service"
)

func main() {
	application := service.Bootstrap()
	server.RunHttpServer(":8000", application)
}
