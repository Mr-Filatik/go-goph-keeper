// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// LoginScreen описывает экран входа и необходимые ему данные.
type LoginScreen struct {
	mainModel     *ViewModel
	LoginInput    textinput.Model
	PasswordInput textinput.Model
	ErrMessage    string
}

// NewLoginScreen создаёт новый экзепляр *LoginScreen.
func NewLoginScreen(mod *ViewModel) *LoginScreen {
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

// String выводит окно и его содержимое в виде строки.
func (s *LoginScreen) String() string {
	view := "\n[Login] Enter email and password:\n"

	view += s.LoginInput.View() + "\n"
	view += s.PasswordInput.View() + "\n\n"

	if s.ErrMessage != "" {
		view += "\n[ERROR]: " + s.ErrMessage + "\n"
	}

	return view
}

// GetHints выводит подсказки по управлению для текущего окна.
func (s *LoginScreen) GetHints() []Hint {
	return []Hint{
		{"Login", []string{KeyEnter}},
		{"Switch", []string{KeyTab, KeyDown, KeyUp}},
		{"Back", []string{KeyEscape}},
		{"Quit", []string{KeyQuit}},
	}
}

// Action описывает логику работы с командами для текущего окна.
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

			// Авторизация пользователя
			screen := s.mainModel.screenLoading

			screen.title = "Login user"
			screen.desc = "Login user by email and password"

			screen.OnProgress = func(percent float64, _ string) tea.Cmd {
				return tea.Tick(RefreshTime, func(time.Time) tea.Msg {
					return loginStep(ctx, percent, login, password)
				})
			}

			screen.OnDone = func(_ any) {
				s.mainModel.currentUser = &user{Login: login}
				s.mainModel.screenPassList.LoadScreen(nil)
			}

			screen.OnCancel = func() {
				s.LoadScreen(func() {
					s.mainModel.screenLogin.ErrMessage = "Login canceled"
					s.mainModel.screenLogin.LoginInput.SetValue(login)
				})
			}

			screen.OnError = func(err error) {
				s.ErrMessage = err.Error()
				s.LoadScreen(func() {
					s.mainModel.screenLogin.ErrMessage = "ERROR" + login + password + err.Error()
				})
			}

			s.mainModel.SetCurrentScreen(screen)

			return s.mainModel, tea.Tick(RefreshTime, func(time.Time) tea.Msg {
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
