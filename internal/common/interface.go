// Package common предоставляет общий функционал для приложений.
package common

import (
	"context"
	"io"
)

// IStarter - интерфейс для всех запускаемых компонентов.
type IStarter interface {
	// Метод для запуска компонентов.
	Start(ctx context.Context) error
}

// IShutdowner - интерфейс для всех останавливаемых компонентов.
type IShutdowner interface {
	// Метод для мягкой остановки компонентов.
	Shutdown(ctx context.Context) error

	// Метод для жёсткой остановки компонентов.
	io.Closer
}
