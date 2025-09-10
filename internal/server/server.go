// Package server предоставляет функционал для запуска приложения сервера.
package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mr-filatik/go-goph-keeper/internal/common"
	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
	"github.com/mr-filatik/go-goph-keeper/internal/server/config"
	"github.com/mr-filatik/go-goph-keeper/internal/server/crypto/jwt"
	"github.com/mr-filatik/go-goph-keeper/internal/server/storage"
)

// IServer - интерфейс для всех серверов приложения.
type IServer interface {
	// Запуск сервера.
	common.IStarter

	// Корректная остановка сервера.
	common.IShutdowner
}

// Run запускает приложение сервера.
func Run() {
	log, logErr := logger.NewZapSugarLogger(logger.LevelDebug, os.Stdout)
	if logErr != nil {
		panic(logErr)
	}

	defer func() {
		if logErr := log.Close(); logErr != nil {
			panic(logErr)
		}
	}()

	exitCtx, exitFn := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer exitFn()

	appConfig := config.Initialize()

	log.Info("Application starting...")

	encr := jwt.NewEncryptor(appConfig.CryptoJWTKey)

	stor := storage.NewMemoryStorage()

	var server IServer

	httpConfig := &HTTPServerConfig{
		Address:   appConfig.ServerAddress,
		Encryptor: encr,
	}

	server = NewHTTPServer(httpConfig, stor, log)

	startErr := server.Start(exitCtx)
	if startErr != nil {
		log.Error("Server starting error", startErr)
	}

	// Ожидание сигнала остановки
	<-exitCtx.Done()
	exitFn()

	log.Info("Application shutdown starting...")
}
