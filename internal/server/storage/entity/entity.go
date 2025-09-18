// Package entity предоставляет сущности.
package entity

import (
	"time"

	"github.com/google/uuid"
)

// User описывает пользователя на сервере.
type User struct {
	ID           string
	Email        string
	PasswordHash string // password hash
}

// NewUser создаёт нового пользователя с уникальным ID.
func NewUser(email, passHash string) *User {
	user := &User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: passHash,
	}

	return user
}

// Token описывает токены для пользователей.
type Token struct{}

// ItemType описывает тип хранимой информации.
type ItemType string

const (
	// ItemLogin - пара логин/пароль.
	ItemLogin ItemType = "login"

	// ItemText - произвольный текст.
	ItemText ItemType = "text"

	// ItemBinary - произвольный бинарь.
	ItemBinary ItemType = "binary"

	// ItemCard - данные банковских карт.
	ItemCard ItemType = "card"
)

// VaultItem описывает тип для хранения информации.
type VaultItem struct {
	ID          string            `json:"id"`
	OwnerID     string            `json:"-"` // владелец (userID)
	Type        ItemType          `json:"type"`
	Title       string            `json:"title"`
	Description string            `json:"desc"`
	Meta        map[string]string `json:"meta"` // произвольная метаинфа
	Username    string            `json:"username"`
	Data        string            `json:"data"` // шифротекст (opaque)
	Version     int64             `json:"version"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}
