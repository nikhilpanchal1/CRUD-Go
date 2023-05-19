package model

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Item struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

// New item with auto gen id
func NewItem(name string, price float64) *Item {
	return &Item{
		ID:    uuid.New().String(), //there is an alternative os/exec that can be used for the same. https://www.geeksforgeeks.org/generate-uuid-in-golang/#
		Name:  name,
		Price: price,
	}
}

// User with auto gen id and hashed passwd
func NewUser(username, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       uuid.New().String(),
		Username: username,
		Password: string(hashedPassword),
	}, nil
}
