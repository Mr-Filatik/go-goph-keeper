// Package handler_test предоставляет функционал для тестирования обработчиков.
package handler_test

import (
	"testing"

	"github.com/mr-filatik/go-goph-keeper/internal/server/handler"
	"github.com/mr-filatik/go-goph-keeper/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewHTTPHandler(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()

	hndlr := handler.NewHTTPHandler(mockLogger)

	assert.NotEmpty(t, hndlr)
}
