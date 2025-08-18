// Package logger предоставляет функционал для логирования.
package logger

import "io"

// LogLevel описывает уровень логирования.
type LogLevel uint32

// Константы - уровни логирования.
const (
	// Уровень логирования debug.
	LevelDebug LogLevel = iota

	// Уровень логирования info.
	LevelInfo

	// Уровень логирования warning.
	LevelWarn

	// Уровень логирования error.
	LevelError
)

// Logger описывает интерфейс для всех логгеров используемых в проекте.
type Logger interface {
	// Debug логирует сообщение и параметры с уровнем debug.
	Debug(message string, keysAndValues ...any)

	// Info логирует сообщение и параметры с уровнем info.
	Info(message string, keysAndValues ...any)

	// Warn логирует сообщение и параметры с уровнем warn и возможной (некритичной) ошибкой.
	Warn(message string, err error, keysAndValues ...any)

	// Error логирует сообщение и параметры с уровнем error и критичной ошибкой.
	Error(message string, err error, keysAndValues ...any)

	// Close (из io.Closer) освобождает используемые логгером ресурсы.
	io.Closer
}

// GetLevelName преобразует уровень логирования в строку.
//
// Параметры:
//   - logLevel: уровень логирования.
func GetLevelName(logLevel LogLevel) string {
	switch logLevel {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warning"
	case LevelError:
		return "error"
	default:
		return "none"
	}
}

// CorrectLevel ограничивает уровень максимально допустимым, если указан неверный.
//
// Параметры:
//   - logLevel: уровень логирования.
func CorrectLevel(logLevel LogLevel) LogLevel {
	switch logLevel {
	case LevelDebug:
		return LevelDebug
	case LevelInfo:
		return LevelInfo
	case LevelWarn:
		return LevelWarn
	case LevelError:
		return LevelError
	default:
		return LevelError
	}
}
