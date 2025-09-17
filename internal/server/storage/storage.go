// Package storage предоставляет функциональность хранилища.
package storage

import (
	"context"
	"errors"
	"time"

	"github.com/mr-filatik/go-goph-keeper/internal/server/storage/entity"
)

// Возможные ошибки при работе с хранилищем.
var (
	ErrEntityAlreadyExists = errors.New("entity already exists")
	ErrEntityNotFound      = errors.New("entity not found")
)

// IUserStorage - интерфейс для всех хранилищ с пользователями.
type IUserStorage interface {
	AddNewUser(ctx context.Context, user *entity.User) (string, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)

	AddNewToken(ctx context.Context, userID string, token *entity.Token) (string, error)
	IsTokenByUserID(ctx context.Context, userID string) bool // add err for database error
	DeleteToken(ctx context.Context, userID string) error
}

// IStorage - интерфейс для всех хранилищ приложения.
type IStorage interface {
	CreateItem(ctx context.Context, it *entity.VaultItem) (string, error)
	UpdateItem(ctx context.Context, it *entity.VaultItem) error
	UpsertItem(ctx context.Context, it *entity.VaultItem) (string, error) // удобно для sync
	GetItem(ctx context.Context, ownerID, id string) (*entity.VaultItem, error)
	ListItems(ctx context.Context, ownerID string) ([]*entity.VaultItem, error)
	DeleteItem(ctx context.Context, ownerID, id string) error

	// ListChangedSince нужен для синхронизации «что изменилось после T»
	ListChangedSince(
		ctx context.Context,
		ownerID string,
		since time.Time,
	) ([]*entity.VaultItem, error)
}
