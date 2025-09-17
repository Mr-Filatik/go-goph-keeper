// Package vault предоставляет функционал для обработчиков запросов для работы с хранилищем.
package vault

import "github.com/mr-filatik/go-goph-keeper/internal/server/storage/entity"

type upsertReq struct {
	ID      string            `json:"id"`
	Type    entity.ItemType   `json:"type"`
	Title   string            `json:"title"`
	Meta    map[string]string `json:"meta"`
	Data    []byte            `json:"data"`
	Version int64             `json:"version"`
}
