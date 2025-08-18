// Package testutil_test предоставляет функционал для тестирования.
package testutil_test

import (
	"testing"

	"github.com/mr-filatik/go-goph-keeper/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewMockWriter проверяет создание *MockWriter.
func TestNewMockWriter(t *testing.T) {
	t.Parallel()

	writer := testutil.NewMockWriter()

	assert.NotNil(t, writer)
}

// TestMockWriter_Write проверяет запись в буфер и слайс.
func TestMockWriter_Write(t *testing.T) {
	t.Parallel()

	writer := testutil.NewMockWriter()
	count := len(writer.Logs)

	number, err := writer.Write([]byte("hello"))

	require.NoError(t, err)
	assert.Equal(t, 5, number)
	require.Len(t, writer.Logs, count+1)

	lastLog := writer.Logs[len(writer.Logs)-1]

	assert.Equal(t, []byte("hello"), lastLog)
}
