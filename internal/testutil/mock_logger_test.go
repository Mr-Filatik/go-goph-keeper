// Package testutil_test предоставляет функционал для тестирования.
package testutil_test

import (
	"errors"
	"testing"

	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
	"github.com/mr-filatik/go-goph-keeper/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errTest = errors.New("test error")

// TestNewMockLogger проверяет создание *MockLogger.
func TestNewMockLogger(t *testing.T) {
	t.Parallel()

	log := testutil.NewMockLogger()

	assert.NotNil(t, log)
}

// TestMockLogger_Debug тестирует функцию *MockLogger.Debug().
func TestMockLogger_Debug(t *testing.T) {
	t.Parallel()

	log := testutil.NewMockLogger()
	count := len(log.Logs)

	log.Debug("message debug", "key debug", "value debug")

	require.Len(t, log.Logs, count+1)

	lastLog := log.Logs[len(log.Logs)-1]

	assert.Equal(t, logger.LevelDebug, lastLog.Level)
	assert.Equal(t, "message debug", lastLog.Message)
	assert.Equal(t, []any{"key debug", "value debug"}, lastLog.Keyvals)
	assert.NoError(t, lastLog.Err)
}

// TestMockLogger_Info тестирует функцию *MockLogger.Info().
func TestMockLogger_Info(t *testing.T) {
	t.Parallel()

	log := testutil.NewMockLogger()
	count := len(log.Logs)

	log.Info("message info", "key info", "value info")

	require.Len(t, log.Logs, count+1)

	lastLog := log.Logs[len(log.Logs)-1]

	assert.Equal(t, logger.LevelInfo, lastLog.Level)
	assert.Equal(t, "message info", lastLog.Message)
	assert.Equal(t, []any{"key info", "value info"}, lastLog.Keyvals)
	assert.NoError(t, lastLog.Err)
}

// TestMockLogger_Warn тестирует функцию *MockLogger.Warn().
func TestMockLogger_Warn(t *testing.T) {
	t.Parallel()

	log := testutil.NewMockLogger()
	count := len(log.Logs)

	log.Warn("message warn", errTest, "key warn", "value warn")

	require.Len(t, log.Logs, count+1)

	lastLog := log.Logs[len(log.Logs)-1]

	assert.Equal(t, logger.LevelWarn, lastLog.Level)
	assert.Equal(t, "message warn", lastLog.Message)
	assert.Equal(t, []any{"key warn", "value warn"}, lastLog.Keyvals)
	require.Error(t, lastLog.Err)
	assert.ErrorIs(t, lastLog.Err, errTest)
}

// TestMockLogger_Error тестирует функцию *MockLogger.Error().
func TestMockLogger_Error(t *testing.T) {
	t.Parallel()

	log := testutil.NewMockLogger()
	count := len(log.Logs)

	log.Error("message error", errTest, "key error", "value error")

	require.Len(t, log.Logs, count+1)

	lastLog := log.Logs[len(log.Logs)-1]

	assert.Equal(t, logger.LevelError, lastLog.Level)
	assert.Equal(t, "message error", lastLog.Message)
	assert.Equal(t, []any{"key error", "value error"}, lastLog.Keyvals)
	require.Error(t, lastLog.Err)
	assert.ErrorIs(t, lastLog.Err, errTest)
}

// TestMockLogger_Close тестирует функцию *MockLogger.Close().
func TestMockLogger_Close(t *testing.T) {
	t.Parallel()

	log := testutil.NewMockLogger()
	closeErr := log.Close()

	assert.NoErrorf(t, closeErr, "Close() error = %v, wantErr %v", closeErr, nil)
}
