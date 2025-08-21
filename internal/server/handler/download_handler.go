// Package handler предоставляет функционал для обработчиков запросов.
package handler

import (
	"errors"
	"net/http"
	"strings"

	root "github.com/mr-filatik/go-goph-keeper"
	"github.com/mr-filatik/go-goph-keeper/internal/common"
)

var errUnsupportedOS = errors.New("unsupported OS")

// ClientInfo отдаёт информацию о поддержимаемых OS.
func (h *HTTPHandler) ClientInfo(writer http.ResponseWriter, _ *http.Request) {
	info := map[string]string{
		"download": "/client/{os}",
		"os":       "windows, macos, linux",
		"example":  "http://localhost:8080/client/linux",
	}
	h.responceWithJSON(writer, info)
}

// ClientDownload отдаёт бинарник CLI-клиента в зависимости от ОС.
func (h *HTTPHandler) ClientDownload(writer http.ResponseWriter, req *http.Request) {
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
		h.responceBadRequest(writer, errUnsupportedOS)

		return
	}

	filepath = strings.Join([]string{root.DirStatic, "/", filename}, "")

	ok, err := common.ExistsFS(root.EmbedStatic, filepath)
	if !ok {
		if err != nil {
			h.responceInternalServerError(writer, err)

			return
		}

		h.responceNotFound(writer, err)

		return
	}

	writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	writer.Header().Set("Content-Type", contentType)
	writer.Header().Set("Content-Transfer-Encoding", "binary")
	writer.Header().Set("Cache-Control", "no-cache")

	http.ServeFile(writer, req, filepath)
}
