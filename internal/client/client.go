// Package client предоставляет функционал для запуска приложения клиента.
package client

import (
	"os"

	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
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

	log.Info("Client starting...")

	log.Info("Client startup completed successfully")
}
