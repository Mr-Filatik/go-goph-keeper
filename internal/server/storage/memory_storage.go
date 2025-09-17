// Package storage предоставляет функциональность хранилища.
package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mr-filatik/go-goph-keeper/internal/server/storage/entity"
)

// MemoryStorage описывает хранилище.
type MemoryStorage struct {
	mu     sync.RWMutex
	users  map[string]*entity.User  // email -> user
	tokens map[string]*entity.Token // userID -> token
	items  map[string]map[string]*entity.VaultItem
}

// NewMemoryStorage создаёт и инициализирует новый экзепляр *MemoryStorage.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		mu:     sync.RWMutex{},
		users:  make(map[string]*entity.User),
		tokens: make(map[string]*entity.Token),
		items:  make(map[string]map[string]*entity.VaultItem),
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
func (m *MemoryStorage) AddNewToken(
	_ context.Context,
	userID string,
	token *entity.Token,
) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.tokens[userID]; ok {
		return "", fmt.Errorf("token: %w", ErrEntityAlreadyExists)
	}

	m.tokens[userID] = token

	return "", nil
}

// IsTokenByUserID производит поиск токена для пользователя по UserID.
func (m *MemoryStorage) IsTokenByUserID(_ context.Context, userID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.tokens[userID]

	return ok
}

// DeleteToken удаляет токен.
func (m *MemoryStorage) DeleteToken(_ context.Context, userID string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	delete(m.tokens, userID)

	return nil
}

// CreateItem создаёт новую запись с паролем.
func (m *MemoryStorage) CreateItem(_ context.Context, item *entity.VaultItem) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.items[item.OwnerID]; !ok {
		m.items[item.OwnerID] = make(map[string]*entity.VaultItem)
	}

	if item.ID == "" {
		item.ID = uuid.New().String()
	}

	if _, ok := m.items[item.OwnerID][item.ID]; ok {
		return "", fmt.Errorf("item: %w", ErrEntityAlreadyExists)
	}

	item.Version = 1
	item.UpdatedAt = time.Now().UTC()
	cl := *item
	m.items[item.OwnerID][item.ID] = &cl

	return item.ID, nil
}

// UpdateItem обновляет текущую запись с паролем.
func (m *MemoryStorage) UpdateItem(_ context.Context, item *entity.VaultItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	userItems := m.items[item.OwnerID]
	if userItems == nil {
		return fmt.Errorf("item: %w", ErrEntityNotFound)
	}

	old := userItems[item.ID]
	if old == nil {
		return fmt.Errorf("item: %w", ErrEntityNotFound)
	}
	// оптимистичная проверка версии:
	// if it.Version != 0 && it.Version != old.Version { return ErrConflict }
	newIt := *item
	newIt.Version = old.Version + 1
	newIt.UpdatedAt = time.Now().UTC()
	userItems[item.ID] = &newIt

	return nil
}

// UpsertItem обновляет или обновляет текущую запись с паролем для синхронизации.
func (m *MemoryStorage) UpsertItem(ctx context.Context, item *entity.VaultItem) (string, error) {
	if item.ID == "" {
		return m.CreateItem(ctx, item)
	}

	if err := m.UpdateItem(ctx, item); err != nil {
		if errors.Is(err, ErrEntityNotFound) {
			return m.CreateItem(ctx, item)
		}

		return "", err
	}

	return item.ID, nil
}

// GetItem получает текущую запись с паролем по ID.
func (m *MemoryStorage) GetItem(
	_ context.Context,
	ownerID, userID string,
) (*entity.VaultItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	userItems := m.items[ownerID]
	if userItems == nil {
		return nil, fmt.Errorf("item: %w", ErrEntityNotFound)
	}

	it := userItems[userID]
	if it == nil {
		return nil, fmt.Errorf("item: %w", ErrEntityNotFound)
	}

	cp := *it

	return &cp, nil
}

// ListItems получает список записей с паролями по пользователю.
func (m *MemoryStorage) ListItems(_ context.Context, ownerID string) ([]*entity.VaultItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	userItems := m.items[ownerID]
	if userItems == nil {
		return []*entity.VaultItem{}, nil
	}

	res := make([]*entity.VaultItem, 0, len(userItems))

	for _, v := range userItems {
		cp := *v
		res = append(res, &cp)
	}

	return res, nil
}

// DeleteItem удаляет текущую запись с паролем по ID.
func (m *MemoryStorage) DeleteItem(_ context.Context, ownerID, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.items[ownerID] == nil {
		return nil
	}

	delete(m.items[ownerID], userID)

	return nil
}

// ListChangedSince обновить данные.
func (m *MemoryStorage) ListChangedSince(
	_ context.Context,
	ownerID string,
	since time.Time,
) ([]*entity.VaultItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	userItems := m.items[ownerID]
	if userItems == nil {
		return []*entity.VaultItem{}, nil
	}

	var res []*entity.VaultItem

	for _, it := range userItems {
		if it.UpdatedAt.After(since) {
			cp := *it
			res = append(res, &cp)
		}
	}

	return res, nil
}
