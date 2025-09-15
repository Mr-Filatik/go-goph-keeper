// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// RefreshTime определяет время для повторной перерисовки интерфейса.
const RefreshTime = 200 * time.Millisecond

type (
	// LoadingProgressMsg описывает команду для обновления прогресса окна загрузки.
	LoadingProgressMsg struct {
		Percent float64
		Status  string
	}

	// LoadingDoneMsg описывает команду для закрытия окна загрузки.
	//
	// Где:
	//   - Payload: что угодно: user, результат операции и т.п.;
	//   - Err: ошибка (при отмене используется отдельный метод, а не context.Canceled).
	LoadingDoneMsg struct {
		Payload any
		Err     error
	}
)

// LoadingScreen описывает экран загрузки и необходимые ему данные.
type LoadingScreen struct {
	mainModel *ViewModel

	// отображение
	title   string
	desc    string
	percent float64
	status  string

	cancel context.CancelFunc

	// фабрика команды: получит ctx с возможностью отмены
	Start func(ctx context.Context) tea.Cmd

	// колбэки вызывающей стороны (всё поведение — тут)
	OnProgress func(percent float64, status string) tea.Cmd
	OnDone     func(payload any) // успех
	OnError    func(err error)   // ошибка
	OnCancel   func()            // отмена (Esc/с)
}

// NewLoadingScreen создаёт новый экзепляр *LoadingScreen.
func NewLoadingScreen(m *ViewModel) *LoadingScreen {
	return &LoadingScreen{
		mainModel: m,
	}
}

// LoadScreen требует указания LoadingParams
func (s *LoadingScreen) LoadScreen(fnc func()) {
	s.mainModel.screenCurrent = s
	// s.mainModel.screenCurrent = s // сделать так и убрать лишний возврат

	s.title = "none"
	s.desc = "none"
	s.percent = 0
	s.status = "Starting…"

	if fnc != nil {
		fnc()
	}
}

// String выводит окно и его содержимое в виде строки.
func (s *LoadingScreen) String() string {
	builder := &strings.Builder{}
	fmt.Fprintf(builder, "\n[%s]\n%s\n", s.title, s.desc)

	if s.percent > 0 {
		fmt.Fprintf(builder, "\n[PERCENT]: %.1f%%\n", s.percent)
	}

	if s.status != "" {
		fmt.Fprintf(builder, "\n[STATUS]: %s\n", s.status)
	}

	return builder.String()
}

// GetHints выводит подсказки по управлению для текущего окна.
func (s *LoadingScreen) GetHints() []Hint {
	return []Hint{
		{"Cancel", []string{KeyEscape}},
	}
}

// Action описывает логику работы с командами для текущего окна.
func (s *LoadingScreen) Action(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgType := msg.(type) {
	case LoadingProgressMsg:
		s.percent = msgType.Percent
		s.status = msgType.Status

		return s.mainModel, s.OnProgress(s.percent, s.status)

	case LoadingDoneMsg:
		// очистим cancel; дальше — отдаём поведение наружу
		s.cancel = nil
		if msgType.Err != nil {
			if s.OnError != nil {
				s.OnError(msgType.Err)
			}
		} else {
			if s.OnDone != nil {
				s.OnDone(msgType.Payload)
			}
		}

		return s.mainModel, nil

	case tea.KeyMsg:
		switch msgType.String() {
		case KeyEscape:
			// отменим и выйдем
			if s.OnCancel != nil {
				s.OnCancel()
			}
		}

		return s.mainModel, nil
	}

	return s.mainModel, nil
}
