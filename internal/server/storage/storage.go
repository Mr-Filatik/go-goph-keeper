// Package storage предоставляет функциональность хранилища.
package storage

import (
	"context"
	"errors"

	"github.com/mr-filatik/go-goph-keeper/internal/server/storage/entity"
)

// Возможные ошибки при работе с хранилищем.
var (
	ErrEntityAlreadyExists = errors.New("entity already exists")
	ErrEntityNotFound      = errors.New("entity not found")
)

// IStorage - интерфейс для всех хранилищ приложения.
type IStorage interface {
	AddNewUser(ctx context.Context, user *entity.User) (string, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)

	AddNewToken(ctx context.Context, userID string, token *entity.Token) (string, error)
	IsTokenByUserID(ctx context.Context, userID string) bool
	DeleteToken(ctx context.Context, userID string) error
}
