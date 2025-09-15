package view

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	title     string
	desc      string

	percent float64
	status  string

	cancel context.CancelFunc

	action tea.Cmd

	// храним логин/пароль, если они нужны для реальной проверки
	login string
	pass  string
}

func NewLoadingScreen(m *teaModel) *LoadingScreen {
	return &LoadingScreen{
		mainModel: m,
	}
}

func (s *LoadingScreen) LoadScreen(fnc func()) IScreen {
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

	return s
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
	// if s.mainModel.loadingCmd != nil {
	// 	cmd := s.mainModel.loadingCmd
	// 	s.mainModel.loadingCmd = nil

	// 	return s.mainModel, cmd
	// }

	switch msg.(type) {
	case LoadingProgressMsg:
		//return s.mainModel, s.action
		// шаг прогресса
		s.percent += 5
		if s.percent >= 100 {
			// тут реальная проверка логина/пароля вместо заглушки:
			if s.login == "demo" && s.pass == "demo" {
				// успех → на список
				s.mainModel.currentUser = &user{Login: s.login}
				s.mainModel.screenCurrent = s.mainModel.screenPassList.LoadScreen(nil)
				s.cancel = nil
				return s.mainModel, nil
			}
			// ошибка → назад на логин с сообщением
			if lscr := s.mainModel.screenLogin; lscr != nil {
				lscr.ErrMessage = "invalid credentials"
				s.mainModel.screenCurrent = lscr.LoadScreen(nil)
			}
			s.cancel = nil
			return s.mainModel, nil
		}

		// обновим статус и запланируем следующий тик
		switch {
		case s.percent < 25:
			s.status = "Authorizing…"
		case s.percent < 50:
			s.status = "Encrypting…"
		case s.percent < 75:
			s.status = "Syncing…"
		default:
			s.status = "Finalizing…"
		}
		return s.mainModel, tea.Tick(150*time.Millisecond, func(time.Time) tea.Msg { return LoadingProgressMsg{} })
	}

	if key, isKey := msg.(tea.KeyMsg); isKey {
		switch key.String() {
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

// func startActionCmd(ctx context.Context) tea.Cmd {
// 	return func() tea.Msg {
// 		ticker := time.NewTicker(150 * time.Millisecond)
// 		defer ticker.Stop()

// 		percent := 0.0
// 		step := 2.5

// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return canceledMsg{}
// 			case <-ticker.C:
// 				percent += step + (rand.Float64()*1.5 - 0.75)
// 				if percent > 100 {
// 					percent = 100
// 				}
// 				// отправим статус раз в ~10%
// 				var st tea.Msg
// 				switch {
// 				case percent < 25:
// 					st = statusMsg{text: "Preparing…"}
// 				case percent < 50:
// 					st = statusMsg{text: "Encrypting…"}
// 				case percent < 75:
// 					st = statusMsg{text: "Uploading…"}
// 				default:
// 					st = statusMsg{text: "Finalizing…"}
// 				}
// 				// вернём пачкой: сначала прогресс, затем статус
// 				// (в цикле Bubble Tea они применятся по очереди)
// 				if percent >= 100 {
// 					return tea.Batch(
// 						func() tea.Msg { return progressMsg{percent: percent} },
// 						func() tea.Msg { return st },
// 						func() tea.Msg { return doneMsg{} },
// 					)()
// 				}
// 				return tea.Batch(
// 					func() tea.Msg { return progressMsg{percent: percent} },
// 					func() tea.Msg { return st },
// 					startActionCmd(ctx), // рекурсивно продолжаем
// 				)()
// 			}
// 		}
// 	}
// }
