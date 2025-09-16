// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mr-filatik/go-goph-keeper/internal/client/service"
)

// LoginScreen описывает экран входа и необходимые ему данные.
type LoginScreen struct {
	mainModel     *MainModel
	LoginInput    textinput.Model
	PasswordInput textinput.Model
	ErrMessage    string
}

// NewLoginScreen создаёт новый экзепляр *LoginScreen.
func NewLoginScreen(mod *MainModel) *LoginScreen {
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

// ValidateScreenData проверяет и корректирует данные для текущего экрана.
func (s *LoginScreen) ValidateScreenData() {}

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
		{"Switch", []string{KeyTab}},
		{"Back", []string{KeyEscape}},
	}
}

// Action описывает логику работы с командами для текущего окна.
func (s *LoginScreen) Action(msg tea.Msg) (*MainModel, tea.Cmd) {
	key, isKey := msg.(tea.KeyMsg)
	if isKey {
		switch key.String() {
		case KeyEscape:
			return s.mainModel.ExitToStartScreen(context.Background())

		case KeyEnter:
			login := s.LoginInput.Value()
			password := s.PasswordInput.Value()

			if login == "" || password == "" {
				s.ErrMessage = "login and password are required"

				return s.mainModel, nil
			}

			ctx, cancelFn := context.WithCancel(context.Background())

			return s.login(ctx, login, password, cancelFn)

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

func (s *LoginScreen) login(
	ctx context.Context,
	login, pass string,
	cancelFn context.CancelFunc,
) (*MainModel, tea.Cmd) {
	prevScreen := s.mainModel.screenLogin

	// Авторизация пользователя
	screen := s.mainModel.screenLoading

	screen.title = "Login user"
	screen.desc = "Login user by email and password"
	screen.percent = 0
	screen.status = "Send request for login..."
	screen.OnProgress = func(percent float64, _ string) tea.Cmd {
		return tea.Tick(RefreshTime, func(time.Time) tea.Msg {
			return s.loginStep(ctx, percent, login, pass)
		})
	}
	screen.OnDone = func(_ any) tea.Cmd {
		s.mainModel.currentUser = &user{Login: login}

		_, cmd := s.loadPasswords(ctx, cancelFn)
		s.mainModel.pendingCmd = cmd

		return nil
	}
	screen.OnCancel = func() {
		cancelFn()

		prevScreen.ErrMessage = "operation canceled"

		s.mainModel.SetCurrentScreen(prevScreen)
	}
	screen.OnError = func(err error) {
		prevScreen.ErrMessage = err.Error()

		s.mainModel.SetCurrentScreen(prevScreen)
	}

	s.mainModel.SetCurrentScreen(screen)

	return s.mainModel, tea.Tick(RefreshTime, func(time.Time) tea.Msg {
		return s.loginStep(ctx, 0, login, pass)
	})
}

func (s *LoginScreen) loginStep(ctx context.Context, _ float64, login, pass string) tea.Msg {
	err := s.mainModel.service.Login(ctx, login, pass)
	if err != nil {
		return LoadingDoneMsg{
			Payload: nil,
			Err:     err,
		}
	}

	return LoadingDoneMsg{
		Payload: nil,
		Err:     nil,
	}
}

func (s *LoginScreen) loadPasswords(
	ctx context.Context,
	cancelFn context.CancelFunc,
) (*MainModel, tea.Cmd) {
	prevScreen := s.mainModel.screenLogin

	// Авторизация пользователя
	screen := s.mainModel.screenLoading

	screen.title = "Load passwords"
	screen.desc = "Load passwords for user " + s.mainModel.currentUser.Login
	screen.percent = 0
	screen.status = "Send request for loading..."
	screen.OnProgress = func(percent float64, _ string) tea.Cmd {
		return tea.Tick(RefreshTime, func(time.Time) tea.Msg {
			return s.loadPasswordsStep(ctx)
		})
	}
	screen.OnDone = func(payload any) tea.Cmd {
		nextScreen := s.mainModel.screenPassList

		items, _ := payload.([]service.Password)
		if items == nil {
			items = []service.Password{}
		}
		nextScreen.Items = items

		s.mainModel.SetCurrentScreen(nextScreen)

		return tea.Tick(0, func(time.Time) tea.Msg { return nil })
	}
	screen.OnCancel = func() {
		cancelFn()

		prevScreen.ErrMessage = "operation canceled"

		s.mainModel.SetCurrentScreen(prevScreen)
	}
	screen.OnError = func(err error) {
		prevScreen.ErrMessage = err.Error()

		s.mainModel.SetCurrentScreen(prevScreen)
	}

	s.mainModel.SetCurrentScreen(screen)

	return s.mainModel, tea.Tick(RefreshTime, func(time.Time) tea.Msg {
		return s.loadPasswordsStep(ctx)
	})
}

func (s *LoginScreen) loadPasswordsStep(ctx context.Context) tea.Msg {
	items, err := s.mainModel.service.GetPasswords(ctx)
	if err != nil {
		return LoadingDoneMsg{
			Payload: []service.Password{},
			Err:     err,
		}
	}

	return LoadingDoneMsg{
		Payload: items,
		Err:     nil,
	}
}
