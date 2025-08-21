package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	root "github.com/mr-filatik/go-goph-keeper"
	"github.com/mr-filatik/go-goph-keeper/internal/common"
)

var errUnsupportedOS = errors.New("unsupported OS")

// ClientInfo отдаёт информацию о поддержимаемых OS.
func (s *HTTPServer) ClientInfo(writer http.ResponseWriter, _ *http.Request) {
	info := map[string]string{
		"download": "/client/{os}",
		"os":       "windows, macos, linux",
		"example":  "http://localhost:8080/client/linux",
	}
	s.serverResponceWithJSON(writer, info)
}

// ClientDownload отдаёт бинарник CLI-клиента в зависимости от ОС.
func (s *HTTPServer) ClientDownload(writer http.ResponseWriter, req *http.Request) {
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
		s.serverResponceBadRequest(writer, errUnsupportedOS)

		return
	}

	filepath = strings.Join([]string{root.DirStatic, "/", filename}, "")

	ok, err := common.ExistsFS(root.EmbedStatic, filepath)
	if !ok {
		if err != nil {
			s.serverResponceInternalServerError(writer, err)

			return
		}

		s.serverResponceNotFound(writer, err)

		return
	}

	writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	writer.Header().Set("Content-Type", contentType)
	writer.Header().Set("Content-Transfer-Encoding", "binary")
	writer.Header().Set("Cache-Control", "no-cache")

	http.ServeFile(writer, req, filepath)
}

func (s *HTTPServer) serverResponceBadRequest(writer http.ResponseWriter, err error) {
	s.log.Error("Bad request error (code 400)", err)
	http.Error(writer, "Error: "+err.Error(), http.StatusBadRequest)
}

func (s *HTTPServer) serverResponceNotFound(writer http.ResponseWriter, err error) {
	s.log.Error("Bad request error (code 404)", err)
	http.Error(writer, "Error: "+err.Error(), http.StatusNotFound)
}

func (s *HTTPServer) serverResponceInternalServerError(writer http.ResponseWriter, err error) {
	s.log.Error("Internal server error (code 500)", err)
	http.Error(writer, "Error: "+err.Error(), http.StatusInternalServerError)
}

func (s *HTTPServer) serverResponceWithJSON(writer http.ResponseWriter, data any) {
	res, err := json.Marshal(data)
	if err != nil {
		s.serverResponceInternalServerError(writer, err)

		return
	}

	writer.WriteHeader(http.StatusOK)

	_, err = writer.Write(res)
	if err != nil {
		s.serverResponceInternalServerError(writer, err)

		return
	}
}
