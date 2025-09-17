// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
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
	mainModel *MainModel

	title   string
	desc    string
	percent float64
	status  string

	OnProgress func(percent float64, status string) tea.Cmd // изменение статуса или прогресса операции
	OnDone     func(payload any)                            // завершение операции с успехом
	OnError    func(err error)                              // завершение операции с ошибкой
	OnCancel   func()                                       // отмена операции
}

// NewLoadingScreen создаёт новый экзепляр *LoadingScreen.
func NewLoadingScreen(m *MainModel) *LoadingScreen {
	return &LoadingScreen{
		mainModel:  m,
		title:      "title",
		desc:       "desc",
		percent:    0,
		status:     "status",
		OnProgress: nil,
		OnDone:     nil,
		OnError:    nil,
		OnCancel:   nil,
	}
}

// ValidateScreenData проверяет и корректирует данные для текущего экрана.
func (s *LoadingScreen) ValidateScreenData() {
	if s.percent < 0 {
		s.percent = 0
	}
}

// String выводит окно и его содержимое в виде строки.
func (s *LoadingScreen) String() string {
	builder := &strings.Builder{}
	fmt.Fprintf(builder, "\n[%s]\n%s\n", s.title, s.desc)

	if s.percent >= 0 {
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

// Update описывает логику работы с командами для текущего окна.
func (s *LoadingScreen) Update(msg tea.Msg) (*MainModel, tea.Cmd) {
	switch msgType := msg.(type) {
	case LoadingProgressMsg:
		s.percent = msgType.Percent
		s.status = msgType.Status

		return s.mainModel, s.OnProgress(s.percent, s.status)

	case LoadingDoneMsg:
		if msgType.Err != nil {
			if s.OnError != nil {
				s.OnError(msgType.Err)

				return s.mainModel, nil
			}
		} else {
			if s.OnDone != nil {
				s.OnDone(msgType.Payload)

				return s.mainModel, nil
			}
		}

		return s.mainModel, nil
		// return s.mainModel, tea.Quit // here

	case tea.KeyMsg:
		if msgType.String() == KeyEscape {
			if s.OnCancel != nil {
				s.OnCancel()
			}
		}

		return s.mainModel, nil
	}

	return s.mainModel, nil
}
