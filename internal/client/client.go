// Package client предоставляет функционал для запуска приложения клиента.
package client

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mr-filatik/go-goph-keeper/internal/common"
	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
	"github.com/mr-filatik/go-goph-keeper/internal/common/repeater"
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

	//appConfig := config.Initialize()

	log.Info("Application starting...",
		"Build Version", buildVersion,
		"Build Date", buildDate,
		"Build Commit", buildCommit,
	)

	attempts := 0
	maxAttempts := 5

	rep := repeater.New[string, string]().
		SetCondition(func(err error) bool {
			return err == nil // повторять, пока не будет nil
		}).
		SetFunc(func(ctx context.Context, s string) (string, error) {
			select {
			case <-ctx.Done():
				return "", ctx.Err() // будет context.DeadlineExceeded при пер-таймауте
			case <-time.After(2 * time.Second): // «задержка» операции
			}

			attempts++
			if attempts < maxAttempts {
				return "", errors.New("temp error")
			}
			return s, nil
		}).
		SetDurationLimit(1*time.Second, 10*time.Second)

	ctx, cancelFn := context.WithCancel(exitCtx)
	defer cancelFn()

	doneCh, retryCh := rep.Run(ctx, "input")

	for {
		select {
		case ev, ok := <-retryCh:
			if ok {
				msg := fmt.Sprintf("Repeat %d, wait %s", ev.Attempt, ev.Wait)
				log.Warn(msg, ev.Err)
			}
		case fin := <-doneCh:
			if fin.Err != nil {
				log.Error("Done with error", fin.Err) // context.DeadlineExceeded / attempts over / ваша ошибка
			} else {
				log.Info("Done with result " + fin.Result)
			}
			return
		}
	}

	// clientConfig := &resty.ClientConfig{
	// 	ServerAddress: appConfig.ServerAddress,
	// }

	// // Add client.
	// mainClient := resty.NewClient(clientConfig, log)

	// startErr := mainClient.Start(exitCtx)
	// if startErr != nil {
	// 	log.Error("Client starting error", startErr)
	// }

	// // Add service with main application logic.
	// mainService := memory.NewService(log)

	// model := view.NewMainModel(mainService)

	// modelErr := model.Start()
	// if modelErr != nil {
	// 	log.Error("View starting error", modelErr)
	// }

	// // Ожидание сигнала остановки
	// <-exitCtx.Done()
	// exitFn()

	// log.Info("Application shutdown starting...")
}
