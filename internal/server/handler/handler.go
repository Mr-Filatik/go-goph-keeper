// Package handler предоставляет общий функционал для обработчиков запросов.
package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
	"github.com/mr-filatik/go-goph-keeper/internal/server/storage"
)

// Handler содержит общие данные для всех хендлеров.
type Handler struct {
	Log  logger.Logger
	Stor storage.IStorage
}

// NewHandler создаёт и инициализирует новый экзепляр *Handler.
//
// Параметры:
//   - log: логгер.
func NewHandler(stor storage.IStorage, log logger.Logger) *Handler {
	return &Handler{
		Log:  log,
		Stor: stor,
	}
}

// ResponseError формирует ответ при ошибках сервера и дополнительно логирует ошибку.
func (h *Handler) ResponseError(writer http.ResponseWriter, code int, err error) {
	msg := fmt.Sprintf("Response error (HTTP code %d) reason: %s", code, err.Error())
	h.Log.Error(msg, err)
	http.Error(writer, "Error", code)
}

// ResponceWithJSON формирует успешный ответ отправляя данные в формате JSON.
func (h *Handler) ResponceWithJSON(writer http.ResponseWriter, data any) {
	res, err := json.Marshal(data)
	if err != nil {
		h.ResponseError(writer, http.StatusInternalServerError, err)

		return
	}

	writer.WriteHeader(http.StatusOK)

	_, err = writer.Write(res)
	if err != nil {
		h.ResponseError(writer, http.StatusInternalServerError, err)

		return
	}
}

// GetDataFromBodyJSON получает данные из запроса в формате JSON.
func GetDataFromBodyJSON[DataType any](req *http.Request, data *DataType) error {
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(req.Body); err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	if err := json.Unmarshal(buf.Bytes(), data); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	return nil
}
