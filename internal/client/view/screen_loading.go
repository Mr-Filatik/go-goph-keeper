package view

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	LoadingProgressMsg struct {
		Percent float64
		Status  string
	}
	LoadingDoneMsg struct {
		Payload any // что угодно: user, результат операции и т.п.
		Err     error
	}
	// canceledMsg struct{}
)

type LoadingScreen struct {
	mainModel *teaModel

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

func NewLoadingScreen(m *teaModel) *LoadingScreen {
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

	// если передан стартовый команд — запустим его из main.Update (см. ниже)
	// if startCmd != nil {
	// 	s.mainModel.loadingCmd = startCmd
	// }
}

func (s *LoadingScreen) String() string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "\n[%s]\n%s\n", s.title, s.desc)
	if s.percent > 0 {
		fmt.Fprintf(b, "\n[PERCENT]: %.1f%%\n", s.percent)
	}
	if s.status != "" {
		fmt.Fprintf(b, "\n[STATUS]: %s\n", s.status)
	}
	fmt.Fprint(b, "\n[Esc]/[c] Cancel   [q] Quit\n")
	return b.String()
}

func (s *LoadingScreen) GetHints() []Hint {
	return []Hint{
		{"Cancel", []string{KeyEscape}},
		{"Quit", []string{KeyQuit}},
	}
}

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
		case KeyQuit:
			// отменим и выйдем
			if s.cancel != nil {
				s.cancel()
			}

			return s.mainModel, tea.Quit

		case KeyEscape, "c":
			if s.cancel != nil {
				s.cancel()
			}

			return s.mainModel, nil
		}
	}

	return s.mainModel, nil
}
