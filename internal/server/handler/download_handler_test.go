// Package handler_test предоставляет функционал для тестирования обработчиков.
package handler_test

import (
	"net/http"
	"testing"

	"github.com/mr-filatik/go-goph-keeper/internal/server/handler"
)

func TestHTTPHandler_ClientInfo(t *testing.T) {
	t.Parallel()

	type args struct {
		writer http.ResponseWriter
		in1    *http.Request
	}

	tests := []struct {
		name string
		h    *handler.HTTPHandler
		args args
	}{}

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			internalTest.h.ClientInfo(internalTest.args.writer, internalTest.args.in1)
		})
	}
}

func TestHTTPHandler_ClientDownload(t *testing.T) {
	t.Parallel()

	type args struct {
		writer http.ResponseWriter
		req    *http.Request
	}

	tests := []struct {
		name string
		h    *handler.HTTPHandler
		args args
	}{}

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			internalTest.h.ClientDownload(internalTest.args.writer, internalTest.args.req)
		})
	}
}
