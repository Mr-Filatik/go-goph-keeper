// Package handler предоставляет общий функционал для обработчиков запросов.
package handler

import (
	"bytes"
	"encoding/json"
	"errors"
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

func (h *Handler) ResponseError(writer http.ResponseWriter, code int, err error) {
	msg := fmt.Sprintf("Response error (HTTP code %d) reason: %s", code, err.Error())
	h.Log.Error(msg, err)
	http.Error(writer, "Error", code)
}

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

func GetDataFromBodyJSON[DataType any](req *http.Request, data *DataType) error {
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(req.Body); err != nil {
		return errors.New(err.Error())
	}

	if err := json.Unmarshal(buf.Bytes(), data); err != nil {
		return errors.New(err.Error())
	}

	return nil
}
