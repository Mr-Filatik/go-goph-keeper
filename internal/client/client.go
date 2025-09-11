// Package client предоставляет функционал для запуска приложения клиента.
package client

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mr-filatik/go-goph-keeper/internal/client/client/http/resty"
	"github.com/mr-filatik/go-goph-keeper/internal/client/config"
	"github.com/mr-filatik/go-goph-keeper/internal/common"
	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
)

//nolint:gochecknoglobals // подстановка линкерных флагов через -ldflags
var (
	buildVersion = "N/A" // Версия сборки приложения.
	buildDate    = "N/A" // Дата сборки приложения.
	buildCommit  = "N/A" // Коммит сборки приложения.
)

// IClient - интерфейс для всех серверов приложения.
type IClient interface {
	// Запуск клиента.
	common.IStarter

	// Корректная остановка клиента.
	common.IShutdowner
}

// Run запускает приложение клиента.
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

	log.Info("Application starting...",
		"Build Version", buildVersion,
		"Build Date", buildDate,
		"Build Commit", buildCommit,
	)

	clientConfig := &resty.ClientConfig{
		ServerAddress: appConfig.ServerAddress,
	}

	mainClient := resty.NewClient(clientConfig, log)

	startErr := mainClient.Start(exitCtx)
	if startErr != nil {
		log.Error("Client starting error", startErr)
	}

	// Ожидание сигнала остановки
	<-exitCtx.Done()
	exitFn()

	log.Info("Application shutdown starting...")
}
