package models

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var users = []User{
	{
		Email:    "email",
		Password: "password",
	},
	{
		Email:    "       ",
		Password: "password",
	},
	{
		Email:    "rqweqweqweqwe",
		Password: "       ",
	},
}

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("../../.env.test"))
	if err != nil {
		log.Fatalf("Cannot load env error: %v", err)
	}
	err = InitDatabase()
	if err != nil {
		log.Fatalf("Cannot init database error: %v", err)
	}
	err = InitRedis()
	if err != nil {
		log.Fatalf("Cannot init redis error: %v", err)
	}
	m.Run()
}
