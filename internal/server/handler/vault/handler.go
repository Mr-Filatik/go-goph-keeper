// Package vault предоставляет функционал для обработчиков запросов для работы с хранилищем.
package vault

import (
	"net/http"
	"time"

	"github.com/mr-filatik/go-goph-keeper/internal/server/handler"
	"github.com/mr-filatik/go-goph-keeper/internal/server/middleware"
	"github.com/mr-filatik/go-goph-keeper/internal/server/storage"
	"github.com/mr-filatik/go-goph-keeper/internal/server/storage/entity"
)

// Handler хранит данные необходимые для обработчиков.
type Handler struct {
	VStor storage.IStorage
	handler.Handler
}

// NewHandler создаёт новый экземпляр Handler.
func NewHandler(h handler.Handler, vStor storage.IStorage) *Handler {
	return &Handler{
		Handler: h,
		VStor:   vStor,
	}
}

// ListItems выводит все пароли пользователя.
func (h *Handler) ListItems(resp http.ResponseWriter, req *http.Request) {
	uid, _ := middleware.GetUserID(req.Context())

	items, err := h.VStor.ListItems(req.Context(), uid)
	if err != nil {
		h.ResponseError(resp, http.StatusInternalServerError, err)

		return
	}

	h.ResponceWithJSON(resp, items)
}

// GetItem выводит конкретный пароль пользователя.
func (h *Handler) GetItem(resp http.ResponseWriter, req *http.Request) {
	uid, _ := middleware.GetUserID(req.Context())

	id := req.PathValue("id")

	item, err := h.VStor.GetItem(req.Context(), uid, id)
	if err != nil {
		h.ResponseError(resp, http.StatusNotFound, err)

		return
	}

	h.ResponceWithJSON(resp, item)
}

// UpsertItem производит обновление пароля.
func (h *Handler) UpsertItem(resp http.ResponseWriter, req *http.Request) {
	uid, _ := middleware.GetUserID(req.Context())

	var upReq upsertReq

	if err := handler.GetDataFromBodyJSON(req, &upReq); err != nil {
		h.ResponseError(resp, http.StatusBadRequest, err)

		return
	}

	item := &entity.VaultItem{
		ID:          upReq.ID,
		OwnerID:     uid,
		Type:        upReq.Type,
		Title:       upReq.Title,
		Meta:        upReq.Meta,
		Data:        upReq.Data,
		Version:     upReq.Version,
		UpdatedAt:   time.Now().UTC(),
		Description: "",
		Username:    "",
	}

	upsertID, err := h.VStor.UpsertItem(req.Context(), item)
	if err != nil {
		h.ResponseError(resp, http.StatusConflict, err)

		return
	}

	h.ResponceWithJSON(resp, map[string]any{"id": upsertID})
}

// DeleteItem производит удаление пароля.
func (h *Handler) DeleteItem(resp http.ResponseWriter, req *http.Request) {
	uid, _ := middleware.GetUserID(req.Context())

	id := req.PathValue("id")

	if err := h.VStor.DeleteItem(req.Context(), uid, id); err != nil {
		h.ResponseError(resp, http.StatusInternalServerError, err)

		return
	}

	resp.WriteHeader(http.StatusNoContent)
}

// SyncSince синхрогнизирует состояние.
func (h *Handler) SyncSince(resp http.ResponseWriter, req *http.Request) {
	uid, _ := middleware.GetUserID(req.Context())

	sinceStr := req.URL.Query().Get("since")

	since := time.Time{}

	if sinceStr != "" {
		if t, err := time.Parse(time.RFC3339, sinceStr); err == nil {
			since = t
		}
	}

	items, err := h.VStor.ListChangedSince(req.Context(), uid, since)
	if err != nil {
		h.ResponseError(resp, http.StatusInternalServerError, err)

		return
	}

	h.ResponceWithJSON(resp, items)
}
