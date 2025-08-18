// Package logger_test предоставляет функционал для тестирования логгеров.
package logger_test

import (
	"io"
	"testing"

	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
	"github.com/mr-filatik/go-goph-keeper/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewZapSugarLogger тестирует функцию NewZapSugarLogger.
//
// Тестируются создание логгера с разными вариантами уровней логирования,
// в том числе и с несуществующим.
func TestNewZapSugarLogger(t *testing.T) {
	t.Parallel()

	type args struct {
		logLevel logger.LogLevel
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "LevelDebug",
			args: args{
				logLevel: logger.LevelDebug,
			},
		},
		{
			name: "LevelInfo",
			args: args{
				logLevel: logger.LevelInfo,
			},
		},
		{
			name: "LevelWarn",
			args: args{
				logLevel: logger.LevelWarn,
			},
		},
		{
			name: "LevelError",
			args: args{
				logLevel: logger.LevelError,
			},
		},
		{
			name: "LevelUnknown",
			args: args{
				logLevel: logger.LogLevel(999),
			},
		},
	}

	writer := testutil.NewMockWriter()

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			_ = newZapSugarLogger(t, internalTest.args.logLevel, writer)
		})
	}
}

// TestZapSugarLogger_Debug тестирует функцию *NewZapSugar.Debug().
//
// Тестирует, будет ли записывать лог через метод для указанного
// минимального уровня логирования логгера.
func TestZapSugarLogger_Debug(t *testing.T) {
	t.Parallel()

	tests := setupArgsForZapSugarLogger(false, false, false)

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			writer := testutil.NewMockWriter()
			log := newZapSugarLogger(t, internalTest.args.logerMinLevel, writer)

			count := len(writer.Logs)

			log.Debug(internalTest.args.msg)

			if internalTest.want != nil {
				count++
			}

			assert.Len(t, writer.Logs, count)

			if internalTest.want != nil {
				assert.Contains(t, string(writer.Logs[len(writer.Logs)-1]), internalTest.want.msg)
			}
		})
	}
}

// TestZapSugarLogger_Info тестирует функцию *NewZapSugar.Info().
//
// Тестирует, будет ли записывать лог через метод для указанного
// минимального уровня логирования логгера.
func TestZapSugarLogger_Info(t *testing.T) {
	t.Parallel()

	tests := setupArgsForZapSugarLogger(true, false, false)

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			writer := testutil.NewMockWriter()
			log := newZapSugarLogger(t, internalTest.args.logerMinLevel, writer)

			count := len(writer.Logs)

			log.Info(internalTest.args.msg)

			if internalTest.want != nil {
				count++
			}

			assert.Len(t, writer.Logs, count)

			if internalTest.want != nil {
				assert.Contains(t, string(writer.Logs[len(writer.Logs)-1]), internalTest.want.msg)
			}
		})
	}
}

// TestZapSugarLogger_Warn тестирует функцию *NewZapSugar.Warn().
//
// Тестирует, будет ли записывать лог через метод для указанного
// минимального уровня логирования логгера.
func TestZapSugarLogger_Warn(t *testing.T) {
	t.Parallel()

	tests := setupArgsForZapSugarLogger(true, true, false)

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			writer := testutil.NewMockWriter()
			log := newZapSugarLogger(t, internalTest.args.logerMinLevel, writer)

			count := len(writer.Logs)

			log.Warn(internalTest.args.msg, nil)

			if internalTest.want != nil {
				count++
			}

			assert.Len(t, writer.Logs, count)

			if internalTest.want != nil {
				assert.Contains(t, string(writer.Logs[len(writer.Logs)-1]), internalTest.want.msg)
			}
		})
	}
}

// TestZapSugarLogger_Error тестирует функцию *NewZapSugar.Error().
//
// Тестирует, будет ли записывать лог через метод для указанного
// минимального уровня логирования логгера.
func TestZapSugarLogger_Error(t *testing.T) {
	t.Parallel()

	tests := setupArgsForZapSugarLogger(true, true, true)

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			writer := testutil.NewMockWriter()
			log := newZapSugarLogger(t, internalTest.args.logerMinLevel, writer)

			count := len(writer.Logs)

			log.Error(internalTest.args.msg, nil)

			if internalTest.want != nil {
				count++
			}

			assert.Len(t, writer.Logs, count)

			if internalTest.want != nil {
				assert.Contains(t, string(writer.Logs[len(writer.Logs)-1]), internalTest.want.msg)
			}
		})
	}
}

// TestZapSugarLogger_Close тестирует функцию *NewZapSugar.Close().
//
// Проверка на отсутствие ошибок при освобождении логгера.
func TestZapSugarLogger_Close(t *testing.T) {
	t.Parallel()

	writer := testutil.NewMockWriter()
	log := newZapSugarLogger(t, logger.LevelDebug, writer)
	closeErr := log.Close()

	assert.NoErrorf(t, closeErr, "Close() error = %v, wantErr %v", closeErr, nil)
}

// testZapSugarLogger описывает сущность теста.
type testZapSugarLogger struct {
	// Имя эксперимента.
	name string

	// Аргументы эксперимента.
	args argZapSugarLogger

	// Ожидаемый результат эксперимента.
	want *wantZapSugarLogger
}

// argZapSugarLogger описывает сущность аргументов теста.
type argZapSugarLogger struct {
	// Сообщение для логирования.
	msg string

	// Уровень логирования логгера.
	logerMinLevel logger.LogLevel

	// Набор других значений.
	keysAndValues []any
}

// wantZapSugarLogger описывает сущность ожидания теста.
type wantZapSugarLogger struct {
	// Ожидаемое сообщение.
	msg string

	// keysAndValues []any
}

// setupArgsForZapSugarLogger создаёт тесты.
//
// Параметры:
//   - wInfo: ожидается ли вывод лога с уровнем info;
//   - wWarn: ожидается ли вывод лога с уровнем warning;
//   - wErr: ожидается ли вывод лога с уровнем error.
func setupArgsForZapSugarLogger(wInfo, wWarn, wErr bool) []testZapSugarLogger {
	wDebug := true // Для логгера с уровнем debug выводятся все сообщения.

	tests := make([]testZapSugarLogger, 0, 4)

	debugTest := setupArgForZapSugarLogger("LevelDebug", logger.LevelDebug, "message debug", wDebug)

	infoTest := setupArgForZapSugarLogger("LevelInfo", logger.LevelInfo, "message info", wInfo)

	warnTest := setupArgForZapSugarLogger("LevelWarn", logger.LevelWarn, "message warn", wWarn)

	errTest := setupArgForZapSugarLogger("LevelError", logger.LevelError, "message error", wErr)

	tests = append(tests,
		debugTest,
		infoTest,
		warnTest,
		errTest,
	)

	return tests
}

// setupArgForZapSugarLogger устанавливает конкретный эксперимент.
//
// Параметры:
//   - name: название эксперимента;
//   - lvl: уровень логирования логгера;
//   - msg: сообщение;
//   - want: ожидается ли что сообщение будет записано логгером.
func setupArgForZapSugarLogger(
	name string,
	lvl logger.LogLevel,
	msg string,
	want bool,
) testZapSugarLogger {
	test := testZapSugarLogger{
		name: name,
		args: argZapSugarLogger{
			logerMinLevel: lvl,
			msg:           msg,
			keysAndValues: []any{"key", "value"},
		},
		want: nil,
	}
	if want {
		test.want = &wantZapSugarLogger{
			msg: msg,
		}
	}

	return test
}

// newZapSugarLogger вспомогательный метод, позволяющий проверить корректность создания логгера.
func newZapSugarLogger(
	t *testing.T,
	logLevel logger.LogLevel,
	writer io.Writer,
) *logger.ZapSugarLogger {
	t.Helper()

	log, err := logger.NewZapSugarLogger(logLevel, writer)

	require.NoErrorf(t, err, "NewZapSugarLogger() error = %v, wantErr %v", err, nil)
	require.NotEmpty(t, log, "NewZapSugarLogger() = %v, want %v", log, nil)

	return log
}
