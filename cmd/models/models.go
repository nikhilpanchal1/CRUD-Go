package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// TODO: gorm:"primaryKey" || and check if it works with uuid. Might be losing performance there
type Item struct {
	gorm.Model
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"ID"` //gorm:"type:uuid;default:uuid_generate_v4()" json:"id
	Name  string    `json:"name"`
	Price float64   `json:"price"`
}

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"ID"` //uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

// New item with auto gen id
func NewItem(name string, price float64) *Item {
	return &Item{
		ID:    uuid.New(), //there is an alternative os/exec that can be used for the same. https://www.geeksforgeeks.org/generate-uuid-in-golang/#
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
		ID:       uuid.New(),
		Username: username,
		Password: string(hashedPassword),
	}, nil
}
