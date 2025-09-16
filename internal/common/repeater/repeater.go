// Package repeater предоставляет реализацию сущности для повторения действий при ошибках, временных сбоях.
package repeater

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	// ErrActionNotSet действие не установлено, нечего запускать.
	ErrActionNotSet = errors.New("action not set")

	// ErrAttemptsOver все попытки исчерпаны.
	ErrAttemptsOver = errors.New("attempts are over")
)

const (
	// DefaultDurationLimitAttempt время ограничения для одной попытки.
	DefaultDurationLimitAttempt = 2 * time.Second

	// DefaultDurationLimitAll время ограничения для всех попыток.
	DefaultDurationLimitAll = 10 * time.Second

	// MaxDurationLimit максимальное время, если не установлены ограничения.
	MaxDurationLimit = 100 * 365 * 24 * time.Hour
)

// RetryEvent — событие перед повтором.
//
// Срабатывает после завершения предыдущей попытки.
type RetryEvent struct {
	// Attempt указывает на номер попытки (нулевая в него не входит).
	Attempt int

	// Err указывает причину повтора.
	//
	// Может быть следующих типов:
	//   - context.DeadlineExceeded: истёкло временное ограничение для работы;
	Err error

	// Wait указывает время ожидания после окончания текущей неудачной попытки.
	Wait time.Duration
}

// DoneEvent — финальный результат.
type DoneEvent[Tout any] struct {
	// Result содержит результат выполнения функции при успехе.
	Result Tout

	// Err указывает возможную ошибку при выполнении.
	//
	// Может быть следующих типов:
	//   - context.Canceled: отменено извне;
	//   - context.DeadlineExceeded: истёкло временное ограничение для работы;
	//   - context.DeadlineExceeded + attempt: истёкло время определённое на одну попытку;
	//   - context.DeadlineExceeded + all: истёкло время определённое на все попытки;
	Err error
}

// Repeater позвляет повторить действие несколько раз, если при его выполнении не выполнилось условие.
type Repeater[Tin any, Tout any] struct {
	action           func(context.Context, Tin) (Tout, error) // основное действие
	condition        func(error) bool                         // условие для выхода из повторителя
	delays           []time.Duration                          // задержки перед повторением
	DurationLimit    time.Duration                            // максимальное время выполнения (одной попытки)
	DurationLimitAll time.Duration                            // максимальное время выполнения (всех попыток)
}

// New создаёт и инициализирует новый объект *Repeater[Tin, Tout].
//
// Параметры:
//   - log: логгер
func New[Tin any, Tout any]() *Repeater[Tin, Tout] {
	return &Repeater[Tin, Tout]{
		delays:           []time.Duration{500 * time.Millisecond, 1 * time.Second, 2 * time.Second},
		condition:        func(err error) bool { return err == nil },
		action:           nil,
		DurationLimit:    DefaultDurationLimitAttempt,
		DurationLimitAll: DefaultDurationLimitAll,
	}
}

// SetFunc устанавливает основное действие для повторения.
//
// Параметры:
//   - f: функция
func (r *Repeater[Tin, Tout]) SetFunc(f func(context.Context, Tin) (Tout, error)) *Repeater[Tin, Tout] {
	r.action = f

	return r
}

// SetCondition устанавливает условие для повторения.
//
// Параметры:
//   - c: функция условие
func (r *Repeater[Tin, Tout]) SetCondition(c func(error) bool) *Repeater[Tin, Tout] {
	r.condition = c

	return r
}

// SetDelays устанавливает время ожидания между повторами.
//
// Параметры:
//   - ds: ожидания
func (r *Repeater[Tin, Tout]) SetDelays(ds []time.Duration) *Repeater[Tin, Tout] {
	if len(ds) > 0 {
		r.delays = ds
	}

	return r
}

// SetDurationLimit устанавливает максимальное время работы.
//
// Параметры:
//   - dLim: лимит времени для одной попытки
//   - dLimAll: лимит времени для всех попыток
func (r *Repeater[Tin, Tout]) SetDurationLimit(dLim, dLimAll time.Duration) *Repeater[Tin, Tout] {
	r.DurationLimit = dLim
	r.DurationLimitAll = dLimAll

	return r
}

// Run запускает повторитель.
//
// Параметры:
//   - data: данные
func (r *Repeater[Tin, Tout]) Run(parent context.Context, data Tin) (<-chan DoneEvent[Tout], <-chan RetryEvent) {
	retryCh := make(chan RetryEvent, len(r.delays))
	doneCh := make(chan DoneEvent[Tout], 1)

	go func() {
		defer close(retryCh)
		defer close(doneCh)

		var zero Tout

		if r.action == nil {
			doneCh <- DoneEvent[Tout]{Result: zero, Err: ErrActionNotSet}

			return
		}

		ctx, cancel := r.withBudget(parent, r.DurationLimitAll)
		if cancel != nil {
			defer cancel()
		}

		// Первая попытка
		done, prevErr := r.attemptAndMaybeFinish(ctx, data, doneCh)
		if done {
			return
		}

		// Повторы
		for attempt, baseWait := range r.delays {
			// Проверка оставшегося бюджета
			rem := remainingUntil(ctx)

			if rem <= 0 {
				doneCh <- DoneEvent[Tout]{Result: zero, Err: context.DeadlineExceeded}

				return
			}

			wait := minDur(baseWait, rem)

			// Сообщаем о предстоящем повторе (с ошибкой предыдущей попытки)
			notifyRetry(retryCh, attempt+1, prevErr, wait)

			// Ждём или отменяемся
			if err := waitOrCancel(ctx, wait); err != nil {
				doneCh <- DoneEvent[Tout]{Result: zero, Err: err}

				return
			}

			// Следующая попытка
			done, prevErr = r.attemptAndMaybeFinish(ctx, data, doneCh)
			if done {
				return
			}
		}

		doneCh <- DoneEvent[Tout]{Result: zero, Err: ErrAttemptsOver}
	}()

	return doneCh, retryCh
}

// withBudget добавляет общий дедлайн (если задан) и возвращает ctx/cancel.
func (r *Repeater[Tin, Tout]) withBudget(
	parent context.Context,
	d time.Duration,
) (context.Context, context.CancelFunc) {
	if d > 0 {
		return context.WithTimeout(parent, d)
	}

	return parent, nil
}

// attemptAndMaybeFinish выполняет одну попытку;
// если условие выполнено — шлёт результат в doneCh и возвращает done=true.
// В противном случае возвращает done=false и ошибку попытки (для передачи в retry-ивент).
func (r *Repeater[Tin, Tout]) attemptAndMaybeFinish(
	ctx context.Context,
	data Tin,
	doneCh chan<- DoneEvent[Tout],
) (bool, error) {
	result, err := r.runAttempt(ctx, data)
	if r.condition(err) {
		doneCh <- DoneEvent[Tout]{Result: result, Err: err}

		return true, nil
	}

	return false, err
}

// notifyRetry шлёт событие о повторе, не блокируясь.
func notifyRetry(ch chan<- RetryEvent, attempt int, err error, wait time.Duration) {
	select {
	case ch <- RetryEvent{Attempt: attempt, Err: err, Wait: wait}:
	default:
	}
}

// waitOrCancel ждёт интервал или возвращает ошибку отмены/дедлайна.
func waitOrCancel(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return fmt.Errorf("cancel: %w", ctx.Err())
	case <-timer.C:
		return nil
	}
}

// runAttempt запускает одну попытку с учётом DurationLimit и общего бюджета.
//
//nolint:ireturn // возвращаем тип-параметр Tout по контракту *Repeater[Tin, Tout]
func (r *Repeater[Tin, Tout]) runAttempt(parent context.Context, data Tin) (Tout, error) {
	var (
		ctx      = parent
		cancelFn context.CancelFunc
	)

	if r.DurationLimit > 0 {
		rem := remainingUntil(parent)

		per := minDur(r.DurationLimit, rem)
		if per <= 0 {
			var zero Tout

			return zero, context.DeadlineExceeded
		}

		ctx, cancelFn = context.WithTimeout(parent, per)
		defer cancelFn()
	}

	return r.action(ctx, data)
}

func remainingUntil(ctx context.Context) time.Duration {
	if dl, ok := ctx.Deadline(); ok {
		return time.Until(dl)
	}

	return MaxDurationLimit
}

func minDur(first, second time.Duration) time.Duration {
	if first <= 0 {
		return second
	}

	if second <= 0 {
		return first
	}

	if first < second {
		return first
	}

	return second
}
