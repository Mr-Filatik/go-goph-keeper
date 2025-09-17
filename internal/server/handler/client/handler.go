// Package client предоставляет функционал для обработчиков запросов для скачивания клиента.
package client

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	root "github.com/mr-filatik/go-goph-keeper"
	"github.com/mr-filatik/go-goph-keeper/internal/server/file"
	"github.com/mr-filatik/go-goph-keeper/internal/server/handler"
)

// Handler хранит данные необходимые для обработчиков.
type Handler struct {
	files file.IFileStorage
	handler.Handler
}

// HandlerOption представляет дополнительные опции для Handler.
type HandlerOption func(*Handler)

// WithCustomFileStorage устанавливает кастомное файловое хранилище.
func WithCustomFileStorage(stor file.IFileStorage) HandlerOption {
	return func(h *Handler) {
		h.files = stor
	}
}

// NewHandler создаёт новый экземпляр Handler.
func NewHandler(hand handler.Handler, opts ...HandlerOption) *Handler {
	clientHandler := &Handler{
		Handler: hand,
		files:   file.NewClientFileStorage(root.EmbedStatic),
	}

	if opts == nil {
		return clientHandler
	}

	for index := range opts {
		opts[index](clientHandler)
	}

	return clientHandler
}

var errUnsupportedOS = errors.New("unsupported OS")

// ClientInfo отдаёт информацию о поддержимаемых OS.
func (h *Handler) ClientInfo(writer http.ResponseWriter, _ *http.Request) {
	resp := InfoResp{
		Path:    "/client/{os}",
		Example: "http://localhost:8080/client/linux",
		OS: []string{
			"linux",
			"macos",
			"windows",
		},
	}

	writer.Header().Set("Content-Type", "application/json")

	h.ResponceWithJSON(writer, resp)
}

// ClientDownload отдаёт бинарник CLI-клиента в зависимости от ОС.
func (h *Handler) ClientDownload(writer http.ResponseWriter, req *http.Request) {
	osParam := req.PathValue("os")

	path, name, err := h.files.GetFileInfo(osParam)
	if err != nil {
		if errors.Is(err, file.ErrUncorrectClientOS) {
			iErr := fmt.Errorf("%s: %w", osParam, err)
			h.ResponseError(writer, http.StatusBadRequest, iErr)

			return
		}

		h.ResponseError(writer, http.StatusInternalServerError, err)

		return
	}

	var contentType string

	switch osParam {
	case file.ClientOSWindows:
		contentType = "application/x-msdownload"
	case file.ClientOSMacOS:
		contentType = "application/x-mach-binary"
	case file.ClientOSLinux:
		contentType = "application/octet-stream"
	default:
		h.ResponseError(writer, http.StatusBadRequest, errUnsupportedOS)

		return
	}

	writer.Header().Set("Content-Disposition", "attachment; filename="+name)
	writer.Header().Set("Content-Type", contentType)
	writer.Header().Set("Content-Transfer-Encoding", "binary")
	writer.Header().Set("Cache-Control", "no-cache")

	// The following code does not work in tests.
	// fullpath := strings.Join([]string{path, name}, "")
	// http.ServeFile(writer, req, fullpath)

	data, dataErr := h.files.GetFileData(path, name)
	if dataErr != nil {
		h.ResponseError(writer, http.StatusInternalServerError, dataErr)

		return
	}

	http.ServeContent(writer, req, name, time.Time{}, bytes.NewReader(data))
}
