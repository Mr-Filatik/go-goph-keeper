// Package storage предоставляет функциональность хранилища.
package storage

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/mr-filatik/go-goph-keeper/internal/server/storage/entity"
)

// MemoryStorage описывает хранилище.
type MemoryStorage struct {
	mu     sync.RWMutex
	users  map[string]*entity.User  // email -> user
	tokens map[string]*entity.Token // userID -> token
}

// NewMemoryStorage создаёт и инициализирует новый экзепляр *MemoryStorage.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		mu:     sync.RWMutex{},
		users:  make(map[string]*entity.User),
		tokens: make(map[string]*entity.Token),
	}
}

// AddNewUser создаёт нового пользователя.
func (m *MemoryStorage) AddNewUser(_ context.Context, user *entity.User) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.users[strings.ToLower(user.Email)]; ok {
		return "", fmt.Errorf("user: %w", ErrEntityAlreadyExists)
	}

	m.users[strings.ToLower(user.Email)] = user

	return user.ID, nil
}

// FindUserByEmail производит поиск пользователя по Email.
func (m *MemoryStorage) FindUserByEmail(_ context.Context, email string) (*entity.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, ok := m.users[strings.ToLower(email)]
	if !ok {
		return nil, fmt.Errorf("user: %w", ErrEntityNotFound)
	}

	return user, nil
}

// AddNewToken регистрирует новый токен.
func (m *MemoryStorage) AddNewToken(_ context.Context, _ string, _ *entity.Token) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return "", nil
}

// IsTokenByUserID производит поиск токена для пользователя по UserID.
func (m *MemoryStorage) IsTokenByUserID(_ context.Context, _ string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return false
}

// DeleteToken удаляет токен.
func (m *MemoryStorage) DeleteToken(_ context.Context, _ string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return nil
}
