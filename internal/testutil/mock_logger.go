// Package testutil предоставляет функционал для тестирования.
package testutil

import "github.com/mr-filatik/go-goph-keeper/internal/common/logger"

// MockLog представляет собой единицу лога для MockLogger.
type MockLog struct {
	// Ошибка.
	Err error

	// Сообщение.
	Message string

	// Остальные значения.
	Keyvals []any

	// Уровень лога.
	Level logger.LogLevel
}

// MockLogger реализация логгера для тестов.
type MockLogger struct {
	// Логи.
	Logs []MockLog
}

var _ logger.Logger = (*MockLogger)(nil)

// Debug логирует сообщение и параметры в уровне debug.
//
// Параметры:
//   - msg: сообщение;
//   - keysAndValues: дополнительные пары ключ-значение.
func (l *MockLogger) Debug(msg string, keysAndValues ...any) {
	l.Logs = append(l.Logs, MockLog{
		Level:   logger.LevelDebug,
		Message: msg,
		Err:     nil,
		Keyvals: keysAndValues,
	})
}

// Info логирует сообщение и параметры в уровне info.
//
// Параметры:
//   - msg: сообщение;
//   - keysAndValues: дополнительные пары ключ-значение.
func (l *MockLogger) Info(msg string, keysAndValues ...any) {
	l.Logs = append(l.Logs, MockLog{
		Level:   logger.LevelInfo,
		Message: msg,
		Err:     nil,
		Keyvals: keysAndValues,
	})
}

// Warn логирует сообщение и параметры в уровне warning.
//
// Параметры:
//   - msg: сообщение;
//   - err: необязательная ошибка;
//   - keysAndValues: дополнительные пары ключ-значение.
func (l *MockLogger) Warn(msg string, err error, keysAndValues ...any) {
	l.Logs = append(l.Logs, MockLog{
		Level:   logger.LevelWarn,
		Message: msg,
		Err:     err,
		Keyvals: keysAndValues,
	})
}

// Error логирует сообщение и параметры в уровне error.
//
// Параметры:
//   - msg: сообщение;
//   - err: ошибка;
//   - keysAndValues: дополнительные пары ключ-значение.
func (l *MockLogger) Error(msg string, err error, keysAndValues ...any) {
	l.Logs = append(l.Logs, MockLog{
		Level:   logger.LevelError,
		Message: msg,
		Err:     err,
		Keyvals: keysAndValues,
	})
}

// Close освобождает используемые логгером ресурсы.
func (l *MockLogger) Close() error {
	return nil
}
