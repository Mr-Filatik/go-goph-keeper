package view

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginScreen struct {
	mainModel     *teaModel
	LoginInput    textinput.Model
	PasswordInput textinput.Model
	ErrMessage    string
}

func NewLoginScreen(mod *teaModel) *LoginScreen {
	// login input
	loginInput := textinput.New()
	loginInput.Placeholder = "your email"
	loginInput.CharLimit = 64
	loginInput.Focus()

	// password inputs
	passInput := textinput.New()
	passInput.Placeholder = "your password"
	passInput.CharLimit = 64
	passInput.EchoMode = textinput.EchoPassword
	passInput.EchoCharacter = '•'

	return &LoginScreen{
		mainModel:     mod,
		LoginInput:    loginInput,
		PasswordInput: passInput,
		ErrMessage:    "",
	}
}

func (s *LoginScreen) LoadScreen(fnc func()) {
	s.mainModel.screenCurrent = s

	// login input
	loginInput := textinput.New()
	loginInput.Placeholder = "your email"
	loginInput.CharLimit = 64
	loginInput.Focus()

	// password inputs
	passInput := textinput.New()
	passInput.Placeholder = "your password"
	passInput.CharLimit = 64
	passInput.EchoMode = textinput.EchoPassword
	passInput.EchoCharacter = '•'

	s.LoginInput = loginInput
	s.PasswordInput = passInput
	s.ErrMessage = ""

	if fnc != nil {
		fnc()
	}
}

func (s *LoginScreen) String() string {
	view := "\n[Login] Enter email and password:\n"

	view += s.LoginInput.View() + "\n"
	view += s.PasswordInput.View() + "\n\n"

	if s.ErrMessage != "" {
		view += "\n[ERROR]: " + s.ErrMessage + "\n"
	}

	return view
}

func (s *LoginScreen) GetHints() []Hint {
	return []Hint{
		{"Login", []string{KeyEnter}},
		{"Switch", []string{KeyTab, KeyDown, KeyUp}},
		{"Back", []string{KeyEscape}},
		{"Quit", []string{KeyQuit}},
	}
}

func (s *LoginScreen) Action(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, isKey := msg.(tea.KeyMsg)
	if isKey {
		switch key.String() {
		case KeyQuit:
			return s.mainModel, tea.Quit

		case KeyEscape:
			s.mainModel.screenStart.LoadScreen(nil)

			return s.mainModel, nil

		case KeyEnter:
			login := s.LoginInput.Value()
			password := s.PasswordInput.Value()

			if login == "" || password == "" {
				s.ErrMessage = "login and password are required"

				return s.mainModel, nil
			}

			ctx, _ := context.WithCancel(context.Background())

			s.mainModel.screenLoading.LoadScreen(func() {
				s.mainModel.screenLoading.title = "Test"
				s.mainModel.screenLoading.desc = "test test test"

				s.mainModel.screenLoading.OnDone = func(payload any) {
					// успех: кладём пользователя и открываем список
					s.mainModel.currentUser = &user{Login: login}
					s.mainModel.screenPassList.LoadScreen(nil)
				}

				s.mainModel.screenLoading.OnError = func(err error) {
					// ошибка: вернуться на логин и показать сообщение
					s.ErrMessage = err.Error()
					s.LoadScreen(func() {
						s.mainModel.screenLogin.ErrMessage = "ERROR" + login + password + err.Error()
					})
				}

				s.mainModel.screenLoading.OnCancel = func() {
					// отмена: вернуться на логин (можно со своим текстом)
					s.ErrMessage = "Canceled"
					s.LoadScreen(nil)
				}

				s.mainModel.screenLoading.OnProgress = func(percent float64, _ string) tea.Cmd {
					return tea.Tick(150*time.Millisecond, func(time.Time) tea.Msg {
						return loginStep(ctx, percent, login, password)
					})
				}
			})

			return s.mainModel, tea.Tick(150*time.Millisecond, func(time.Time) tea.Msg {
				return loginStep(ctx, 0, login, password)
			})

		case KeyTab, KeyDown, KeyUp:
			if s.LoginInput.Focused() {
				s.LoginInput.Blur()
				s.PasswordInput.Focus()
			} else {
				s.PasswordInput.Blur()
				s.LoginInput.Focus()
			}

			return s.mainModel, nil
		}
	}

	var cmd tea.Cmd
	if s.LoginInput.Focused() {
		s.LoginInput, cmd = s.LoginInput.Update(key)
	} else {
		s.PasswordInput, cmd = s.PasswordInput.Update(key)
	}

	return s.mainModel, cmd
}

func loginStep(ctx context.Context, current float64, login, pass string) tea.Msg {
	select {
	case <-ctx.Done():
		return LoadingDoneMsg{Err: context.Canceled}
	default:
	}

	// увеличиваем прогресс и шлём его
	next := current + 5
	if next >= 100 {
		// финал: делайте реальную проверку логина/пароля
		if login == "demo" && pass == "demo" {
			return LoadingDoneMsg{Payload: user{Login: login}}
		}
		return LoadingDoneMsg{Err: fmt.Errorf("invalid credentials %s %s", login, pass)}
	}

	// промежуточный прогресс — отправляем И планируем следующий тик
	status := "Authorizing…"
	switch {
	case next < 30:
		status = "Authorizing…"
	case next < 60:
		status = "Encrypting…"
	case next < 90:
		status = "Syncing…"
	default:
		status = "Finalizing…"
	}

	return LoadingProgressMsg{
		Percent: next,
		Status:  status,
	}
}
