// Package server предоставляет функционал для запуска приложения сервера.
package server

import (
	"os"

	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
)

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

	log.Info("Server starting...")

	log.Info("Server startup completed successfully")
}
