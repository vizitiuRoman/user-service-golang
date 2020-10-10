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
	Validate(string) error
	Create() (*User, error)
	DeleteByID(uint64) error
	FindAll() (*[]User, error)
	FindByID(uint64) (*User, error)
	FindByEmail(string) (*User, error)
	Update(uint64) error
}

type User struct {
	ID        uint64    `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (user *User) Prepare() {
	user.ID = 0
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.Password = html.EscapeString(strings.TrimSpace(user.Password))
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
	}
	return nil
}

func (user *User) Create() (*User, error) {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return &User{}, err
	}
	user.Password = string(hashedPassword)
	rows, err := db.NamedQuery(`
		INSERT INTO users (email, password)
		VALUES (:email, :password)
		RETURNING email, id;
		`,
		&user,
	)
	if err != nil {
		return &User{}, err
	}
	if rows.Next() {
		rows.StructScan(&user)
	}
	return user, nil
}

func (user *User) FindAll() (*[]User, error) {
	var users []User
	err := db.Select(&users, "SELECT email, id FROM users")
	if err != nil {
		return &users, err
	}
	return &users, nil
}

func (user *User) FindByID(id uint64) (*User, error) {
	err := db.Get(user, "SELECT email, id FROM users WHERE id = $1", id)
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) FindByEmail(email string) (*User, error) {
	err := db.Get(user, "SELECT email, id, password FROM users WHERE email = $1", email)
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) Update(userID uint64) error {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	_, err = db.Query(`
		UPDATE users 
		SET email=$2, password=$3 
		WHERE id = $1`,
		userID, user.Email, user.Password,
	)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) DeleteByID(userID uint64) error {
	_, err := db.Query(`DELETE FROM users WHERE id = $1`,
		userID,
	)
	if err != nil {
		return err
	}
	return nil
}
