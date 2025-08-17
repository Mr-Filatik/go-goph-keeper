package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// ZapSugarLogger реализация логгера через *zap.SugaredLogger.
type ZapSugarLogger struct {
	log      *zap.SugaredLogger // логгер
	logLevel LogLevel           // уровень логирования
}

var _ Logger = (*ZapSugarLogger)(nil)

// NewZapSugarLogger инициализирует и создаёт новый экземпляр *ZapSugarLogger.
//
// Параметры:
//   - logLevel: уровень логирования.
func NewZapSugarLogger(logLevel LogLevel) (*ZapSugarLogger, error) {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("create zap logger error: %w", err)
	}

	zslog := &ZapSugarLogger{
		logLevel: CorrectLevel(logLevel),
		log:      zapLog.Sugar(),
	}
	zslog.Info(
		"ZapSugarLogger creation completed successfully",
		"level", GetLevelName(zslog.logLevel),
	)

	return zslog, nil
}

// Debug логирует сообщение и параметры в уровне debug.
//
// Параметры:
//   - msg: сообщение;
//   - keysAndValues: дополнительные пары ключ-значение.
func (l *ZapSugarLogger) Debug(msg string, keysAndValues ...any) {
	if LevelDebug >= l.logLevel {
		l.log.Infow(msg, keysAndValues...)
	}
}

// Info логирует сообщение и параметры в уровне info.
//
// Параметры:
//   - msg: сообщение;
//   - keysAndValues: дополнительные пары ключ-значение.
func (l *ZapSugarLogger) Info(msg string, keysAndValues ...any) {
	if LevelInfo >= l.logLevel {
		l.log.Infow(msg, keysAndValues...)
	}
}

// Warn логирует сообщение и параметры в уровне warning.
//
// Параметры:
//   - msg: сообщение;
//   - err: необязательная ошибка;
//   - keysAndValues: дополнительные пары ключ-значение.
func (l *ZapSugarLogger) Warn(msg string, err error, keysAndValues ...any) {
	if LevelWarn >= l.logLevel {
		if err != nil {
			keysAndValues = append(keysAndValues, []any{"error", err.Error()})
		}

		l.log.Infow(msg, keysAndValues...)
	}
}

// Error логирует сообщение и параметры в уровне error.
//
// Параметры:
//   - msg: сообщение;
//   - err: ошибка;
//   - keysAndValues: дополнительные пары ключ-значение.
func (l *ZapSugarLogger) Error(msg string, err error, keysAndValues ...any) {
	if LevelError >= l.logLevel {
		if err != nil {
			keysAndValues = append(keysAndValues, []any{"error", err.Error()})
		} else {
			keysAndValues = append(keysAndValues, []any{"error", "nil"})
		}

		l.log.Infow(msg, keysAndValues...)
	}
}

// Close освобождает используемые логгером ресурсы.
func (l *ZapSugarLogger) Close() error {
	if err := l.log.Sync(); err != nil {
		// На Windows есть ошибка при вызове Sync, поэтому игнорирую ошибку.
		// Все ресурсы должны освободиться в любом случае.
		_ = err
	}

	return nil
}
