package vault_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mr-filatik/go-goph-keeper/internal/server/handler"
	"github.com/mr-filatik/go-goph-keeper/internal/server/handler/vault"
	"github.com/mr-filatik/go-goph-keeper/internal/server/middleware"
	"github.com/mr-filatik/go-goph-keeper/internal/server/storage"
	"github.com/mr-filatik/go-goph-keeper/internal/server/storage/entity"
	"github.com/mr-filatik/go-goph-keeper/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
	===== NewHandler =====
*/

func TestVault_NewHandler(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()
	mainHandler := handler.NewHandler(nil, mockLogger)

	h := vault.NewHandler(*mainHandler, nil)
	assert.NotNil(t, h)
}

/*
	===== Helpers =====
*/

func withUser(req *http.Request) *http.Request {
	return req.WithContext(middleware.WithUserID(req.Context(), "user-1"))
}

func mustJSONBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()

	var buf bytes.Buffer

	require.NoError(t, json.NewEncoder(&buf).Encode(v))

	return &buf
}

func assertHTTP(t *testing.T, gotStatus int, gotBody string) {
	t.Helper()

	assert.Equal(t, http.StatusOK, gotStatus)

	assert.NotEmpty(t, gotBody)
}

/*
	===== Mocks =====
*/

type mockStorage struct {
	CreateItemFn       func(ctx context.Context, it *entity.VaultItem) (string, error)
	UpdateItemFn       func(ctx context.Context, it *entity.VaultItem) error
	UpsertItemFn       func(ctx context.Context, it *entity.VaultItem) (string, error)
	GetItemFn          func(ctx context.Context, ownerID, id string) (*entity.VaultItem, error)
	ListItemsFn        func(ctx context.Context, ownerID string) ([]*entity.VaultItem, error)
	DeleteItemFn       func(ctx context.Context, ownerID, id string) error
	ListChangedSinceFn func(ctx context.Context, ownerID string, since time.Time) ([]*entity.VaultItem, error)
}

func (m *mockStorage) CreateItem(ctx context.Context, it *entity.VaultItem) (string, error) {
	return m.CreateItemFn(ctx, it)
}

func (m *mockStorage) UpdateItem(ctx context.Context, it *entity.VaultItem) error {
	return m.UpdateItemFn(ctx, it)
}

func (m *mockStorage) UpsertItem(ctx context.Context, it *entity.VaultItem) (string, error) {
	return m.UpsertItemFn(ctx, it)
}

func (m *mockStorage) GetItem(ctx context.Context, ownerID, id string) (*entity.VaultItem, error) {
	return m.GetItemFn(ctx, ownerID, id)
}

func (m *mockStorage) ListItems(ctx context.Context, ownerID string) ([]*entity.VaultItem, error) {
	return m.ListItemsFn(ctx, ownerID)
}

func (m *mockStorage) DeleteItem(ctx context.Context, ownerID, id string) error {
	return m.DeleteItemFn(ctx, ownerID, id)
}

func (m *mockStorage) ListChangedSince(
	ctx context.Context,
	ownerID string,
	since time.Time,
) ([]*entity.VaultItem, error) {
	return m.ListChangedSinceFn(ctx, ownerID, since)
}

/*
	===== Handler.ListItems =====
*/

func TestVault_ListItems(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()
	mainHandler := handler.NewHandler(nil, mockLogger)

	mockStor := &mockStorage{
		ListItemsFn: func(_ context.Context, ownerID string) ([]*entity.VaultItem, error) {
			return []*entity.VaultItem{
				{ID: "1", OwnerID: ownerID, Title: "Email", Type: entity.ItemLogin},
			}, nil
		},
	}
	vaultHandler := vault.NewHandler(*mainHandler, mockStor)

	req := httptest.NewRequest(http.MethodGet, "/vault/items", http.NoBody)
	req = withUser(req)
	rr := httptest.NewRecorder()

	vaultHandler.ListItems(rr, req)

	assertHTTP(t, rr.Code, rr.Body.String())
}

func TestVault_ListItems_StorageError(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()
	mainHandler := handler.NewHandler(nil, mockLogger)

	mockStor := &mockStorage{
		ListItemsFn: func(_ context.Context, _ string) ([]*entity.VaultItem, error) {
			return nil, assert.AnError
		},
	}
	vaultHandler := vault.NewHandler(*mainHandler, mockStor)

	req := httptest.NewRequest(http.MethodGet, "/vault/items", http.NoBody)
	req = withUser(req)
	rr := httptest.NewRecorder()

	vaultHandler.ListItems(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "Error\n", rr.Body.String())
}

/*
	===== Handler.GetItem =====
*/

func TestVault_GetItem(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()
	mainHandler := handler.NewHandler(nil, mockLogger)

	mockStor := &mockStorage{
		GetItemFn: func(_ context.Context, ownerID, id string) (*entity.VaultItem, error) {
			//nolint:exhaustruct // not all fields needed in test
			item := &entity.VaultItem{
				ID:      id,
				OwnerID: ownerID,
				Title:   "Email",
				Type:    entity.ItemLogin,
			}

			return item, nil
		},
	}
	vaultHandler := vault.NewHandler(*mainHandler, mockStor)

	req := httptest.NewRequest(http.MethodGet, "/vault/items/1", http.NoBody)

	ctx, ok := any(middleware.WithUserID(context.Background(), "user-1")).(context.Context)
	require.True(t, ok)

	req = req.WithContext(context.WithValue(req.Context(), ctx, nil)) // защитный no-op

	req = withUser(req)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	vaultHandler.GetItem(rr, req)

	assertHTTP(t, rr.Code, rr.Body.String())
}

func TestVault_GetItem_NotFound(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()
	mainHandler := handler.NewHandler(nil, mockLogger)

	mockStor := &mockStorage{
		GetItemFn: func(_ context.Context, _, _ string) (*entity.VaultItem, error) {
			return nil, storage.ErrEntityNotFound
		},
	}

	vaultHandler := vault.NewHandler(*mainHandler, mockStor)

	req := httptest.NewRequest(http.MethodGet, "/vault/items/404", http.NoBody)
	req = withUser(req)
	req.SetPathValue("id", "404")

	rr := httptest.NewRecorder()
	vaultHandler.GetItem(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "Error\n", rr.Body.String())
}

/*
	===== Handler.UpsertItem =====
*/

func TestVault_UpsertItem_OK(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()
	mainHandler := handler.NewHandler(nil, mockLogger)

	mockStor := &mockStorage{
		UpsertItemFn: func(_ context.Context, it *entity.VaultItem) (string, error) {
			require.Equal(t, "user-1", it.OwnerID)

			return "new-id", nil
		},
	}

	vaultHandler := vault.NewHandler(*mainHandler, mockStor)

	body := map[string]any{
		"id":      "",
		"type":    "login",
		"title":   "Email",
		"meta":    map[string]string{"site": "example.com"},
		"data":    []byte{0x01, 0x02},
		"version": 0,
	}
	req := httptest.NewRequest(http.MethodPost, "/vault/items", mustJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req = withUser(req)

	rr := httptest.NewRecorder()
	vaultHandler.UpsertItem(rr, req)

	assertHTTP(t, rr.Code, rr.Body.String())
}

func TestVault_UpsertItem_BadRequest(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()
	mainHandler := handler.NewHandler(nil, mockLogger)
	vaultHandler := vault.NewHandler(*mainHandler, &mockStorage{
		UpsertItemFn: func(_ context.Context, _ *entity.VaultItem) (string, error) {
			return "", nil
		},
	})

	// битый JSON
	req := httptest.NewRequest(http.MethodPost, "/vault/items", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	req = withUser(req)

	rr := httptest.NewRecorder()
	vaultHandler.UpsertItem(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "Error\n", rr.Body.String())
}

func TestVault_UpsertItem_Conflict(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()
	mainHandler := handler.NewHandler(nil, mockLogger)
	vaultHandler := vault.NewHandler(*mainHandler, &mockStorage{
		UpsertItemFn: func(_ context.Context, _ *entity.VaultItem) (string, error) {
			return "", storage.ErrEntityAlreadyExists
		},
	})

	body := map[string]any{"id": "1", "type": "login", "title": "Email"}
	req := httptest.NewRequest(http.MethodPost, "/vault/items", mustJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req = withUser(req)

	rr := httptest.NewRecorder()
	vaultHandler.UpsertItem(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, "Error\n", rr.Body.String())
}

/*
	===== Handler.DeleteItem =====
*/

func TestVault_DeleteItem_OK(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()
	mainHandler := handler.NewHandler(nil, mockLogger)
	vaultHandler := vault.NewHandler(*mainHandler, &mockStorage{
		DeleteItemFn: func(_ context.Context, _, _ string) error { return nil },
	})

	req := httptest.NewRequest(http.MethodDelete, "/vault/items/1", http.NoBody)
	req = withUser(req)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	vaultHandler.DeleteItem(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body.String())
}

func TestVault_DeleteItem_Error(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()
	mainHandler := handler.NewHandler(nil, mockLogger)
	vaultHandler := vault.NewHandler(*mainHandler, &mockStorage{
		DeleteItemFn: func(_ context.Context, _, _ string) error { return assert.AnError },
	})

	req := httptest.NewRequest(http.MethodDelete, "/vault/items/1", http.NoBody)
	req = withUser(req)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	vaultHandler.DeleteItem(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "Error\n", rr.Body.String())
}

/*
	===== Handler.SyncSince =====
*/

func TestVault_SyncSince_OK(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()
	mainHandler := handler.NewHandler(nil, mockLogger)
	now := time.Now().UTC()

	vaultHandler := vault.NewHandler(*mainHandler, &mockStorage{
		ListChangedSinceFn: func(_ context.Context, ownerID string, since time.Time) ([]*entity.VaultItem, error) {
			assert.Equal(t, "user-1", ownerID)
			assert.False(t, since.IsZero())

			return []*entity.VaultItem{{ID: "1", OwnerID: ownerID, UpdatedAt: now}}, nil
		},
	})

	path := "/vault/sync?since=" + now.Add(-time.Minute).Format(time.RFC3339)
	req := httptest.NewRequest(http.MethodGet, path, http.NoBody)
	req = withUser(req)

	rr := httptest.NewRecorder()
	vaultHandler.SyncSince(rr, req)

	assertHTTP(t, rr.Code, rr.Body.String())
}

func TestVault_SyncSince_Error(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()
	mainHandler := handler.NewHandler(nil, mockLogger)

	vaultHandler := vault.NewHandler(*mainHandler, &mockStorage{
		ListChangedSinceFn: func(_ context.Context, _ string, _ time.Time) ([]*entity.VaultItem, error) {
			return nil, assert.AnError
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/vault/sync", http.NoBody)
	req = withUser(req)

	rr := httptest.NewRecorder()
	vaultHandler.SyncSince(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "Error\n", rr.Body.String())
}
