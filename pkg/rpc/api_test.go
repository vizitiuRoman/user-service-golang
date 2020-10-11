package rpc

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	. "github.com/user-service/pkg/models"
)

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

func TestGetUsersRPC(t *testing.T) {
	user := User{Email: "email@gmail.com", Password: "password"}
	_, err := user.Create()
	if err != nil {
		t.Errorf("Cannot create user error: %v", err)
	}

	var users []User
	var userRPC UserRPC
	err = userRPC.GetUsers(int(0), &users)
	if err != nil {
		t.Errorf("Cannot get user error: %v", err)
	}
	assert.Equal(t, len(users), 1)
	err = user.DeleteByID(user.ID)
	if err != nil {
		t.Errorf("Cannot delete user error: %v", err)
	}
}

func TestGetUserRPC(t *testing.T) {
	user := User{Email: "email@gmail.com", Password: "password"}
	createdUser, err := user.Create()
	if err != nil {
		t.Errorf("Cannot create user error: %v", err)
	}

	var foundUser User
	var userRPC UserRPC
	err = userRPC.GetUser(int(createdUser.ID), &foundUser)
	if err != nil {
		t.Errorf("Cannot get user error: %v", err)
	}
	assert.Equal(t, createdUser.Email, foundUser.Email)
	assert.Equal(t, createdUser.ID, foundUser.ID)
	err = user.DeleteByID(createdUser.ID)
	if err != nil {
		t.Errorf("Cannot delete user error: %v", err)
	}
}
