package client

import (
	"errors"
	"net/http"
	"strings"

	root "github.com/mr-filatik/go-goph-keeper"
	"github.com/mr-filatik/go-goph-keeper/internal/common"
	"github.com/mr-filatik/go-goph-keeper/internal/server/handler"
)

// ClientHandler хранит данные необходимые для обработчиков.
type ClientHandler struct {
	handler.Handler
}

// ClientHandlerOption представляет дополнительные опции для Handler.
type ClientHandlerOption func(*ClientHandler)

// func WithLogger(l *slog.Logger) Option { return func(h *Handler){ h.log = l } }

// NewHandler создаёт новый экземпляр Handler.
func NewHandler(hand handler.Handler, opts ...ClientHandlerOption) *ClientHandler {
	h := &ClientHandler{
		Handler: hand,
	}

	for _, o := range opts {
		o(h)
	}

	return h
}

var errUnsupportedOS = errors.New("unsupported OS")

// ClientInfo отдаёт информацию о поддержимаемых OS.
func (h *ClientHandler) ClientInfo(writer http.ResponseWriter, _ *http.Request) {
	info := map[string]string{
		"download": "/client/{os}",
		"os":       "windows, macos, linux",
		"example":  "http://localhost:8080/client/linux",
	}
	h.ResponceWithJSON(writer, info)
}

// ClientDownload отдаёт бинарник CLI-клиента в зависимости от ОС.
func (h *ClientHandler) ClientDownload(writer http.ResponseWriter, req *http.Request) {
	os := req.PathValue("os")

	var filename, filepath, contentType string

	switch os {
	case "windows":
		filename = "client-windows.exe"
		contentType = "application/x-msdownload"
	case "macos":
		filename = "client-macos.exe"
		contentType = "application/x-mach-binary"
	case "linux":
		filename = "client-linux.exe"
		contentType = "application/octet-stream"
	default:
		h.ResponseError(writer, http.StatusBadRequest, errUnsupportedOS)

		return
	}

	filepath = strings.Join([]string{root.DirStatic, "/", filename}, "")

	ok, err := common.ExistsFS(root.EmbedStatic, filepath)
	if !ok {
		if err != nil {
			h.ResponseError(writer, http.StatusInternalServerError, err)

			return
		}

		h.ResponseError(writer, http.StatusNotFound, err)

		return
	}

	writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	writer.Header().Set("Content-Type", contentType)
	writer.Header().Set("Content-Transfer-Encoding", "binary")
	writer.Header().Set("Cache-Control", "no-cache")

	http.ServeFile(writer, req, filepath)
}
