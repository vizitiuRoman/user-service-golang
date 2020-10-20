package rpc

import (
	. "github.com/user-service/pkg/models"
	"go.uber.org/zap"
)

type Rpc interface {
	GetUsers(int, *[]User) error
	GetUser(int, *User) error
}

type UserRPC struct {
	logger *zap.SugaredLogger
}

func NewUserRPC(logger *zap.SugaredLogger) *UserRPC {
	return &UserRPC{logger}
}

func (u *UserRPC) GetUsers(_ int, reply *[]User) error {
	var user User
	users, err := user.FindAll()
	if err != nil {
		return err
	}
	*reply = *users
	return nil
}

func (u *UserRPC) GetUser(userID int, reply *User) error {
	var user User
	foundUser, err := user.FindByID(uint64(userID))
	if err != nil {
		return err
	}
	*reply = *foundUser
	return nil
}
