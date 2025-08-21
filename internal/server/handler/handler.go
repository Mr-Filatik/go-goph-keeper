// Package handler предоставляет функционал для обработчиков запросов.
package handler

import (
	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
)

// HTTPHandler — общие данные хендлеров.
type HTTPHandler struct {
	log logger.Logger
}

// NewHTTPHandler создаёт и инициализирует новый экзепляр *HTTPHandler.
//
// Параметры:
//   - log: логгер.
func NewHTTPHandler(log logger.Logger) *HTTPHandler {
	return &HTTPHandler{
		log: log,
	}
}
