package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyHashedPassword(t *testing.T) {
	password := "password"
	hashedPassword, err := hashPassword(password)
	if err != nil {
		t.Errorf("Cannot hash password error: %v", err)
	}
	err = VerifyPassword(string(hashedPassword), password)
	assert.Nil(t, err)
	err = VerifyPassword(string(hashedPassword), "r")
	assert.NotNil(t, err)
}

func TestPrepareUser(t *testing.T) {
	user := User{
		Email:    "        ",
		Password: "password",
	}
	user.Prepare()
	assert.Equal(t, user.Email, "")
	assert.Equal(t, user.Password, "password")
}

func TestValidateUser(t *testing.T) {
	for _, user := range users {
		user.Prepare()
		err := user.Validate("")
		if user.Password == "" {
			assert.Equal(t, err.Error(), "Required Password")
		}
		if user.Email == "" {
			assert.Equal(t, err.Error(), "Required Email")
		}
		if user.Email != "" && user.Password != "" {
			assert.Nil(t, err)
		}
	}
}

func TestCreateUser(t *testing.T) {

}

func TestDeleteByIDUser(t *testing.T) {

}

func TestFindAllUsers(t *testing.T) {

}

func TestFindUserByID(t *testing.T) {

}

func TestFindUserByEmail(t *testing.T) {

}

func TestUpdateUser(t *testing.T) {

}
