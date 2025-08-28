// Package handler предоставляет функционал для обработчиков запросов.
package handler

import (
	"encoding/json"
	"net/http"
)

func (h *HTTPHandler) responceBadRequest(writer http.ResponseWriter, err error) {
	h.log.Error("Bad request error (code 400)", err)
	http.Error(writer, "Error: "+err.Error(), http.StatusBadRequest)
}

func (h *HTTPHandler) responceNotFound(writer http.ResponseWriter, err error) {
	h.log.Error("Bad request error (code 404)", err)
	http.Error(writer, "Error: "+err.Error(), http.StatusNotFound)
}

func (h *HTTPHandler) responceInternalServerError(writer http.ResponseWriter, err error) {
	h.log.Error("Internal server error (code 500)", err)
	http.Error(writer, "Error: "+err.Error(), http.StatusInternalServerError)
}

func (h *HTTPHandler) responceWithJSON(writer http.ResponseWriter, data any) {
	res, err := json.Marshal(data)
	if err != nil {
		h.responceInternalServerError(writer, err)

		return
	}

	writer.WriteHeader(http.StatusOK)

	_, err = writer.Write(res)
	if err != nil {
		h.responceInternalServerError(writer, err)

		return
	}
}
