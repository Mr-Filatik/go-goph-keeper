package logger

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapSugarLogger реализация логгера через *zap.SugaredLogger.
type ZapSugarLogger struct {
	// Логгер.
	log *zap.SugaredLogger

	// Уровень логирования.
	logLevel LogLevel
}

var _ Logger = (*ZapSugarLogger)(nil)

// NewZapSugarLogger инициализирует и создаёт новый экземпляр *ZapSugarLogger.
//
// Параметры:
//   - logLevel: уровень логирования.
func NewZapSugarLogger(logLevel LogLevel, out io.Writer) (*ZapSugarLogger, error) {
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(out),
		zapcore.DebugLevel,
	)

	zapLog := zap.New(core, zap.AddCaller())

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
		l.log.Debugw(msg, keysAndValues...)
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
			keysAndValues = append(keysAndValues, "error", err.Error())
		}

		l.log.Warnw(msg, keysAndValues...)
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
			keysAndValues = append(keysAndValues, "error", err.Error())
		} else {
			keysAndValues = append(keysAndValues, "error", "nil")
		}

		l.log.Errorw(msg, keysAndValues...)
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
