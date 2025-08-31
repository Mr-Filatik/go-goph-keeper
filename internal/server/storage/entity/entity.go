// Package entity предоставляет сущности.
package entity

import "github.com/google/uuid"

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
type Token struct {
}
