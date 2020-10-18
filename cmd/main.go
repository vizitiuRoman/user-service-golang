package main

import (
	"github.com/user-service/pkg/server"
)

func main() {
	srv := server.NewService()
	srv.Init()
	go srv.StartRPC()
	srv.StartAPI()
}
