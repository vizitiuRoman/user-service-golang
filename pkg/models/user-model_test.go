package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func clearDatabase(t *testing.T) {
	_, err := db.Query("DELETE FROM users")
	if err != nil {
		t.Errorf("Cannot delete * from users error: %v", err)
	}
}

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
	defer clearDatabase(t)

	user := users[0]
	createdUser, err := user.Create()
	if err != nil {
		t.Errorf("Cannot create user error: %v", err)
	}
	assert.Equal(t, user.Email, createdUser.Email)
	assert.NotEmpty(t, createdUser.ID)
	assert.NotEmpty(t, createdUser.Password)
}

func TestDeleteUserByID(t *testing.T) {
	defer clearDatabase(t)

	user := users[0]
	createdUser, err := user.Create()
	if err != nil {
		t.Errorf("Cannot create user error: %v", err)
	}
	err = user.DeleteByID(createdUser.ID)
	assert.Nil(t, err)
}

func TestFindAllUsers(t *testing.T) {
	defer clearDatabase(t)

	var user User
	var samples = []User{
		{
			Email:    "email@gmail.com",
			Password: "password",
		},
		{
			Email:    "email2@gmail.com",
			Password: "password",
		},
		{
			Email:    "email3@gmail.com",
			Password: "password",
		},
		{
			Email:    "email4@gmail.com",
			Password: "password",
		},
	}

	for _, sample := range samples {
		_, err := sample.Create()
		if err != nil {
			t.Errorf("Cannot create user error: %v", err)
		}
	}
	users, err := user.FindAll()
	if err != nil {
		t.Errorf("Cannot find all users error: %v", err)
	}
	assert.Equal(t, len(*users), 4)
}

func TestFindUserByID(t *testing.T) {
	defer clearDatabase(t)

	user := users[0]
	createdUser, err := user.Create()
	if err != nil {
		t.Errorf("Cannot create user error: %v", err)
	}

	foundUser, err := user.FindByID(createdUser.ID)
	if err != nil {
		t.Errorf("Cannot find user error: %v", err)
	}
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, createdUser.ID, foundUser.ID)
}

func TestFindUserByEmail(t *testing.T) {
	defer clearDatabase(t)

	user := users[0]
	createdUser, err := user.Create()
	if err != nil {
		t.Errorf("Cannot create user error: %v", err)
	}

	foundUser, err := user.FindByEmail(createdUser.Email)
	if err != nil {
		t.Errorf("Cannot find user error: %v", err)
	}
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, createdUser.ID, foundUser.ID)
}

func TestUpdateUser(t *testing.T) {
	defer clearDatabase(t)

	user := users[0]
	createdUser, err := user.Create()
	if err != nil {
		t.Errorf("Cannot create user error: %v", err)
	}

	err = createdUser.Update()
	if err != nil {
		t.Errorf("Cannot update user error: %v", err)
	}
}
