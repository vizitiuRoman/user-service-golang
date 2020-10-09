package main

import (
	"github.com/user-service/pkg/server"
)

func main() {
	srv := server.NewService()
	go srv.StartRPC()
	srv.StartAPI()
}
