package main

import (
	"fmt"
	"net/rpc"
	"os"

	"github.com/joho/godotenv"
	. "github.com/user-service/pkg/models"
)

func main() {
	_ = godotenv.Load()

	client, _ := rpc.DialHTTP("tcp", "127.0.0.1:"+os.Getenv("RPC_PORT"))

	var user User
	var users []User

	if err := client.Call("UserRPC.GetUser", 1, &user); err != nil {
		fmt.Println("Error: UserRPC.GetUser", err)
	}
	fmt.Println(user)

	if err := client.Call("UserRPC.GetUsers", 1, &users); err != nil {
		fmt.Println("Error: UserRPC.GetUsers", err)
	}
	fmt.Println(users)
}
