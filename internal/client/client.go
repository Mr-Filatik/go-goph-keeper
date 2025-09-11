// Package client предоставляет функционал для запуска приложения клиента.
package client

import (
	"os"

	"github.com/mr-filatik/go-goph-keeper/internal/client/config"
	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
)

//nolint:gochecknoglobals // подстановка линкерных флагов через -ldflags
var (
	buildVersion = "N/A" // Версия сборки приложения.
	buildDate    = "N/A" // Дата сборки приложения.
	buildCommit  = "N/A" // Коммит сборки приложения.
)

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

	appConfig := config.Initialize()

	log.Info("Application starting...",
		"Build Version", buildVersion,
		"Build Date", buildDate,
		"Build Commit", buildCommit,
	)

	log.Info("Client starting...")

	_ = appConfig.ServerAddress

	log.Info("Client startup completed successfully")
}
