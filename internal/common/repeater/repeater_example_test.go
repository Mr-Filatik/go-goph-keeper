package repeater_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
	"github.com/mr-filatik/go-goph-keeper/internal/common/repeater"
)

var errTemp = errors.New("temp error")

// ExampleRepeater_Run — пример использования *repeater.Repeater.
func ExampleRepeater_Run() {
	log, logErr := logger.NewZapSugarLogger(logger.LevelDebug, os.Stdout)
	if logErr != nil {
		panic(logErr)
	}

	defer func() {
		if logErr := log.Close(); logErr != nil {
			panic(logErr)
		}
	}()

	rep := initRepeater()

	ctx, cancelFn := context.WithCancel(context.Background())
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
				log.Error("Done with error", fin.Err)
			} else {
				log.Info("Done with result " + fin.Result)
			}

			return
		}
	}
}

func initRepeater() *repeater.Repeater[string, string] {
	attempts := 0
	maxAttempts := 3

	return repeater.New[string, string]().
		SetCondition(func(err error) bool {
			return err == nil
		}).
		SetFunc(func(ctx context.Context, data string) (string, error) {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(2 * time.Second):
			}

			attempts++
			if attempts < maxAttempts {
				return "", errTemp
			}

			return data, nil
		}).
		SetDurationLimit(1*time.Second, 10*time.Second)
}
