package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserModel interface {
	Prepare()
	Validate(action string) error
	Create() (*User, error)
	Update() (*User, error)
	DeleteByID(userID uint64) (*User, error)
	FindByEmail(email string) (*User, error)
}

type User struct {
	ID        uint64    `db:"id" json:"-"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password,omitempty"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Age       string    `db:"age" json:"age"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}

func (user *User) Prepare() {
	user.ID = 0
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.Password = html.EscapeString(strings.TrimSpace(user.Password))
	user.LastName = html.EscapeString(strings.TrimSpace(user.LastName))
	user.FirstName = html.EscapeString(strings.TrimSpace(user.FirstName))
}

func (user *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if user.Email == "" {
			return errors.New("Required Email")
		}
		if user.Password == "" {
			return errors.New("Required Password")
		}
	default:
		if user.Email == "" {
			return errors.New("Required Email")
		}
		if user.Password == "" {
			return errors.New("Required Password")
		}
		if user.Age == "" {
			return errors.New("Required Age")
		}
		if user.FirstName == "" {
			return errors.New("Required FirstName")
		}
		if user.LastName == "" {
			return errors.New("Required LastName")
		}
	}
	return nil
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (user *User) Create() (*User, error) {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return &User{}, err
	}
	user.Password = string(hashedPassword)
	_, err = db.NamedExec(`
		INSERT INTO users (email, password, first_name, last_name, age)
		VALUES (:email, :password, :first_name, :last_name, :age)`,
		&user,
	)
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) Update() (*User, error) {
	return user, nil
}

func (user *User) FindByEmail(email string) (*User, error) {
	err := db.Get(user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) DeleteByID(userID uint64) (*User, error) {
	return user, nil
}
