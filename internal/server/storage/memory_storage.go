package storage

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/mr-filatik/go-goph-keeper/internal/server/storage/entity"
)

type MemoryStorage struct {
	mu     sync.RWMutex
	users  map[string]*entity.User  // email -> user
	tokens map[string]*entity.Token // userID -> token
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		mu:     sync.RWMutex{},
		users:  make(map[string]*entity.User),
		tokens: make(map[string]*entity.Token),
	}
}

func (m *MemoryStorage) AddNewUser(ctx context.Context, user *entity.User) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.users[strings.ToLower(user.Email)]; ok {
		return "", fmt.Errorf("user: %w", ErrEntityAlreadyExists)
	}

	user.ID = uuid.New().String()
	m.users[strings.ToLower(user.Email)] = user

	return user.ID, nil
}

func (m *MemoryStorage) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, ok := m.users[strings.ToLower(email)]
	if !ok {
		return nil, fmt.Errorf("user: %w", ErrEntityNotFound)
	}

	return user, nil
}

func (m *MemoryStorage) AddNewToken(ctx context.Context, userID string, token *entity.Token) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return "", nil
}

func (m *MemoryStorage) IsTokenByUserID(ctx context.Context, userID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return false
}

func (m *MemoryStorage) DeleteToken(ctx context.Context, userID string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return nil
}
