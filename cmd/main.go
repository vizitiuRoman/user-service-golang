package main

import (
	"github.com/user-service/pkg/server"
)

func main() {
	srv := server.NewServer()
	srv.StartServer()
}
