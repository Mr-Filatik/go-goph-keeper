// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"context"
	"fmt"

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
	step          int
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
		step:          0,
	}
}

// ValidateScreenData проверяет и корректирует данные для текущего экрана.
func (s *LoginScreen) ValidateScreenData() {
	s.step = 0
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
		{"Switch", []string{KeyTab}},
		{"Back", []string{KeyEscape}},
	}
}

// Update описывает логику работы с командами для текущего окна.
func (s *LoginScreen) Update(msg tea.Msg) (*MainModel, tea.Cmd) {
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
			screen := s.mainModel.screenLoading

			screen.title = "Login user"
			screen.desc = "Login user by email and password"
			screen.percent = 0
			screen.status = "Send request for login..."
			screen.OnProgress = func(percent float64, _ string) tea.Cmd {
				return s.actionCmd(ctx, login, password)
				// return tea.Tick(0, func(time.Time) tea.Msg {
				// 	return s.actionCmd(ctx, cancelFn, login, password) // s.loginStep(ctx, percent, login, password)
				// })
			}
			screen.OnDone = func(payload any) {
				s.mainModel.currentUser = &user{Login: login}

				next := s.mainModel.screenPassList

				items, _ := payload.([]service.Password)
				if items != nil {
					next.Items = items
				} else {
					next.Items = []service.Password{}
				}

				s.mainModel.SetCurrentScreen(next)
				// return tea.Tick(0, func(time.Time) tea.Msg {
				// 	return s.actionCmd(ctx, cancelFn, login, password) // s.loginStep(ctx, percent, login, password)
				// })
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

			// screen.action = s.actionCmd(ctx, cancelFn, login, password)

			return s.mainModel, s.actionCmd(ctx, login, password)

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

func (s *LoginScreen) actionCmd(
	ctx context.Context,
	login string,
	pass string,
) tea.Cmd {
	return func() tea.Msg {
		switch s.step {
		case 0:
			// шаг 1: авторизация
			if err := s.mainModel.service.Login(ctx, login, pass); err != nil {
				return LoadingDoneMsg{Err: fmt.Errorf("авторизация: %w", err)}
			}
			s.step = 1
			return LoadingProgressMsg{Percent: 0.3, Status: "Загрузка данных…"}

		case 1:
			// шаг 2: загрузка паролей
			items, err := s.mainModel.service.GetPasswords(ctx)
			if err != nil {
				return LoadingDoneMsg{Err: fmt.Errorf("загрузка данных: %w", err)}
			}
			// успех
			//s.step = 2
			// return LoadingProgressMsg{Percent: 0.8, Status: "Обработка загруженных данных…"}
			return LoadingDoneMsg{Payload: items, Err: nil}
		}
		// запасной случай
		return LoadingDoneMsg{Payload: nil, Err: nil}
	}
}
