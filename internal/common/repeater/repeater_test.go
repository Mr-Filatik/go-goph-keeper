package repeater_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mr-filatik/go-goph-keeper/internal/common/repeater"
)

var errTemp = errors.New("temp error")

// helper: ждём финальный результат и параллельно собираем retry-события.
func wait(
	doneCh <-chan repeater.DoneEvent[string],
	retryCh <-chan repeater.RetryEvent,
	timeout time.Duration,
) (repeater.DoneEvent[string], []repeater.RetryEvent, error) {
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()

	var events []repeater.RetryEvent

	for {
		select {
		case ev, ok := <-retryCh:
			if ok {
				events = append(events, ev)
			}
		case fin := <-doneCh:
			return fin, events, nil
		case <-deadline.C:
			return repeater.DoneEvent[string]{}, events, context.DeadlineExceeded
		}
	}
}

func TestRun_ActionNotSet(t *testing.T) {
	t.Parallel()

	rep := repeater.New[int, string]() // action не задан

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	t.Cleanup(cancel)

	doneCh, retryCh := rep.Run(ctx, 0)

	fin, events, err := wait(doneCh, retryCh, 500*time.Millisecond)
	if err != nil {
		t.Fatalf("wait failed: %v", err)
	}

	if fin.Err == nil || !errors.Is(fin.Err, repeater.ErrActionNotSet) {
		t.Fatalf("expected ErrActionNotSet, got %v", fin.Err)
	}

	if len(events) != 0 {
		t.Fatalf("expected no retry events, got %d", len(events))
	}
}

func TestRun_SuccessFirstAttempt_NoRetries(t *testing.T) {
	t.Parallel()

	rep := repeater.New[string, string]().
		SetFunc(func(_ context.Context, in string) (string, error) {
			return in + "-ok", nil
		})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	doneCh, retryCh := rep.Run(ctx, "x")

	fin, events, err := wait(doneCh, retryCh, time.Second)
	if err != nil {
		t.Fatalf("wait failed: %v", err)
	}

	if fin.Err != nil {
		t.Fatalf("unexpected error: %v", fin.Err)
	}

	if fin.Result != "x-ok" {
		t.Fatalf("unexpected result: %q", fin.Result)
	}

	if len(events) != 0 {
		t.Fatalf("expected 0 retries, got %d", len(events))
	}
}

func TestRun_RetriesThenSuccess(t *testing.T) {
	t.Parallel()

	attempts := 0
	rep := repeater.New[string, string]().
		SetDelays([]time.Duration{30 * time.Millisecond, 30 * time.Millisecond}).
		SetFunc(func(_ context.Context, input string) (string, error) {
			attempts++
			if attempts < 3 {
				// имитируем быстрый провал
				return "", errTemp
			}

			return input + "-ok", nil
		})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	doneCh, retryCh := rep.Run(ctx, "data")

	fin, events, err := wait(doneCh, retryCh, time.Second)
	if err != nil {
		t.Fatalf("wait failed: %v", err)
	}

	if fin.Err != nil || fin.Result != "data-ok" {
		t.Fatalf("expected success after retries; fin=%+v", fin)
	}
	// Должно быть 2 ретрая (между 3 попытками)
	if len(events) != 2 {
		t.Fatalf("expected 2 retry events, got %d", len(events))
	}
	// Проверим, что Attempt нумеруется с 1 и содержит предыдущую ошибку
	for attempt, event := range events {
		wantAttempt := attempt + 1
		if event.Attempt != wantAttempt {
			t.Fatalf("event #%d: expected Attempt=%d, got %d", attempt, wantAttempt, event.Attempt)
		}

		if event.Err == nil {
			t.Fatalf("event #%d: expected non-nil Err", attempt)
		}
	}
}

func TestRun_PerAttemptTimeoutTriggers(t *testing.T) {
	t.Parallel()

	// Первая попытка должна упасть по таймауту попытки
	rep := repeater.New[string, string]().
		SetDurationLimit(50*time.Millisecond, 2*time.Second). // 50ms на одну попытку
		SetDelays([]time.Duration{10 * time.Millisecond}).
		SetFunc(func(ctx context.Context, input string) (string, error) {
			// ctx-aware "sleep"
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(200 * time.Millisecond): // дольше таймаута попытки
			}

			return input, nil
		})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(cancel)

	doneCh, retryCh := rep.Run(ctx, "x")

	fin, events, err := wait(doneCh, retryCh, 2*time.Second)
	if err != nil {
		t.Fatalf("wait failed: %v", err)
	}

	// Должен быть хотя бы один ретрай с context.DeadlineExceeded как Err в событии
	if len(events) == 0 {
		t.Fatalf("expected at least 1 retry event")
	}

	foundTimeout := false

	for _, event := range events {
		if errors.Is(event.Err, context.DeadlineExceeded) {
			foundTimeout = true

			break
		}
	}

	if !foundTimeout {
		t.Fatalf("expected per-attempt timeout in retry events, got: %+v", events)
	}
	// Итог может быть attempts over или успешный, это зависит от delays/бюджета,
	// нас интересовал именно признак таймаута попытки.
	_ = fin
}

func TestRun_AttemptsOver(t *testing.T) {
	t.Parallel()

	// Небольшие задержки, несколько повторов, всё время ошибка
	rep := repeater.New[int, string]().
		SetDelays([]time.Duration{10 * time.Millisecond, 10 * time.Millisecond}).
		SetFunc(func(_ context.Context, _ int) (string, error) {
			return "in", errTemp
		})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	doneCh, retryCh := rep.Run(ctx, 42)

	fin, events, err := wait(doneCh, retryCh, time.Second)
	if err != nil {
		t.Fatalf("wait failed: %v", err)
	}

	if !errors.Is(fin.Err, repeater.ErrAttemptsOver) {
		t.Fatalf("expected ErrAttemptsOver, got %v", fin.Err)
	}
	// Ретраев должно быть ровно len(delays)
	if got, want := len(events), 2; got != want {
		t.Fatalf("expected %d retry events, got %d", want, got)
	}
}

func TestRun_CanceledFromOutside(t *testing.T) {
	t.Parallel()

	rep := repeater.New[string, string]().
		SetDelays([]time.Duration{50 * time.Millisecond, 50 * time.Millisecond}).
		SetFunc(func(ctx context.Context, _ string) (string, error) {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(80 * time.Millisecond):
			}

			return "", errTemp
		})

	ctx, cancel := context.WithCancel(context.Background())
	// отменяем чуть позже старта
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()

	doneCh, retryCh := rep.Run(ctx, "x")

	fin, events, err := wait(doneCh, retryCh, time.Second)
	if err != nil {
		t.Fatalf("wait failed: %v", err)
	}

	if !errors.Is(fin.Err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", fin.Err)
	}

	_ = events
}

// attempts := 0
// maxAttempts := 5

// rep := repeater.New[string, string]().
// 	SetCondition(func(err error) bool {
// 		return err == nil // повторять, пока не будет nil
// 	}).
// 	SetFunc(func(ctx context.Context, s string) (string, error) {
// 		select {
// 		case <-ctx.Done():
// 			return "", ctx.Err() // будет context.DeadlineExceeded при пер-таймауте
// 		case <-time.After(2 * time.Second): // «задержка» операции
// 		}

// 		attempts++
// 		if attempts < maxAttempts {
// 			return "", errors.New("temp error")
// 		}
// 		return s, nil
// 	}).
// 	SetDurationLimit(1*time.Second, 10*time.Second)

// ctx, cancelFn := context.WithCancel(exitCtx)
// defer cancelFn()

// doneCh, retryCh := rep.Run(ctx, "input")

// for {
// 	select {
// 	case ev, ok := <-retryCh:
// 		if ok {
// 			msg := fmt.Sprintf("Repeat %d, wait %s", ev.Attempt, ev.Wait)
// 			log.Warn(msg, ev.Err)
// 		}
// 	case fin := <-doneCh:
// 		if fin.Err != nil {
// 			log.Error("Done with error", fin.Err) // context.DeadlineExceeded / attempts over / ваша ошибка
// 		} else {
// 			log.Info("Done with result " + fin.Result)
// 		}
// 		return
// 	}
// }
