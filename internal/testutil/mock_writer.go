// Package testutil предоставляет функционал для тестирования.
package testutil

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

// MockWriter реализация oi.Writer для тестов.
type MockWriter struct {
	// Мьютекс.
	mu sync.Mutex

	// Буффер для хранения данных.
	buffer *bytes.Buffer

	// Логи в виде слайса []byte.
	Logs [][]byte
}

var _ io.Writer = (*MockWriter)(nil)

// NewMockWriter инициализирует и создаёт новый экземпляр *MockWriter.
func NewMockWriter() *MockWriter {
	return &MockWriter{
		mu:     sync.Mutex{},
		buffer: &bytes.Buffer{},
		Logs:   make([][]byte, 0),
	}
}

// Write записывает данные в буфер и слайс.
func (w *MockWriter) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	num, err := w.buffer.Write(data)
	if err != nil {
		return num, fmt.Errorf("buffer write error: %w", err)
	}

	w.Logs = append(w.Logs, bytes.TrimSpace(data))

	return num, nil
}
