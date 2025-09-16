// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// RegisterScreen описывает экран рагистрации и необходимые ему данные.
type RegisterScreen struct {
	mainModel     *MainModel
	LoginInput    textinput.Model
	PasswordInput textinput.Model
	ErrMessage    string
}

// NewRegisterScreen создаёт новый экзепляр *RegisterScreen.
func NewRegisterScreen(mod *MainModel) *RegisterScreen {
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

	return &RegisterScreen{
		mainModel:     mod,
		LoginInput:    loginInput,
		PasswordInput: passInput,
		ErrMessage:    "",
	}
}

// ValidateScreenData проверяет и корректирует данные для текущего экрана.
func (s *RegisterScreen) ValidateScreenData() {}

// String выводит окно и его содержимое в виде строки.
func (s *RegisterScreen) String() string {
	view := "\n[Register] Enter email and password:\n"

	view += s.LoginInput.View() + "\n"
	view += s.PasswordInput.View() + "\n\n"

	if s.ErrMessage != "" {
		view += "\n[ERROR]: " + s.ErrMessage + "\n"
	}

	return view
}

// GetHints выводит подсказки по управлению для текущего окна.
func (s *RegisterScreen) GetHints() []Hint {
	return []Hint{
		{"Login", []string{KeyEnter}},
		{"Switch", []string{KeyTab}},
		{"Back", []string{KeyEscape}},
	}
}

// Action описывает логику работы с командами для текущего окна.
func (s *RegisterScreen) Action(msg tea.Msg) (*MainModel, tea.Cmd) {
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

			prevScreen := s.mainModel.screenLogin

			// Авторизация пользователя
			screen := s.mainModel.screenLoading

			screen.title = "Register user"
			screen.desc = "Register user by email and password"
			screen.percent = 0
			screen.status = "Send request for register..."
			screen.OnProgress = func(percent float64, _ string) tea.Cmd {
				return tea.Tick(RefreshTime, func(time.Time) tea.Msg {
					return s.registerStep(ctx, percent, login, password)
				})
			}
			screen.OnDone = func(_ any) tea.Cmd {
				s.mainModel.currentUser = &user{Login: login}
				s.mainModel.SetCurrentScreen(s.mainModel.screenPassList)
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
				return s.registerStep(ctx, 0, login, password)
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

func (s *RegisterScreen) registerStep(ctx context.Context, _ float64, login, pass string) tea.Msg {
	err := s.mainModel.service.Register(ctx, login, pass)
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
