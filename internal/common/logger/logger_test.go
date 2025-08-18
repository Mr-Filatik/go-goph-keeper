// Package logger_test предоставляет функционал для тестирования логгеров.
package logger_test

import (
	"testing"

	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
	"github.com/stretchr/testify/assert"
)

// TestGetLevelName тестирует функцию GetLevelName.
//
// Тестируются разные варианты уровней логирования, в том числе и несуществующий.
func TestGetLevelName(t *testing.T) {
	t.Parallel()

	type args struct {
		logLevel logger.LogLevel
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "LevelDebug",
			args: args{
				logLevel: logger.LevelDebug,
			},
			want: "debug",
		},
		{
			name: "LevelInfo",
			args: args{
				logLevel: logger.LevelInfo,
			},
			want: "info",
		},
		{
			name: "LevelWarn",
			args: args{
				logLevel: logger.LevelWarn,
			},
			want: "warning",
		},
		{
			name: "LevelError",
			args: args{
				logLevel: logger.LevelError,
			},
			want: "error",
		},
		{
			name: "LevelUnknown",
			args: args{
				logLevel: logger.LogLevel(999),
			},
			want: "none",
		},
	}

	for ti := range tests {
		tt := tests[ti]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := logger.GetLevelName(tt.args.logLevel)
			assert.Equalf(t, tt.want, got, "GetLevelName() = %v, want %v", got, tt.want)
		})
	}
}

// TestCorrectLevel тестирует функцию CorrectLevel.
//
// Тестируется корректировка по значению, ограничение по максимально строгому уровню.
func TestCorrectLevel(t *testing.T) {
	t.Parallel()

	type args struct {
		logLevel logger.LogLevel
	}

	tests := []struct {
		name string
		args args
		want logger.LogLevel
	}{
		{
			name: "LevelDebug",
			args: args{
				logLevel: logger.LevelDebug,
			},
			want: logger.LevelDebug,
		},
		{
			name: "LevelInfo",
			args: args{
				logLevel: logger.LevelInfo,
			},
			want: logger.LevelInfo,
		},
		{
			name: "LevelWarn",
			args: args{
				logLevel: logger.LevelWarn,
			},
			want: logger.LevelWarn,
		},
		{
			name: "LevelError",
			args: args{
				logLevel: logger.LevelError,
			},
			want: logger.LevelError,
		},
		{
			name: "LevelUnknown",
			args: args{
				logLevel: logger.LogLevel(999),
			},
			want: logger.LevelError,
		},
	}

	for ti := range tests {
		tt := tests[ti]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := logger.CorrectLevel(tt.args.logLevel)
			assert.Equalf(t, tt.want, got, "CorrectLevel() = %v, want %v", got, tt.want)
		})
	}
}
