// Package view содержит логику для работы с пользовательским интерфейсом.
//
//nolint:dupl // сейчас выглядит как дупликат screen_login, но набор полей должен быть разным.
package view

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mr-filatik/go-goph-keeper/internal/client/service"
)

// RegisterScreen описывает экран рагистрации и необходимые ему данные.
type RegisterScreen struct {
	mainModel     *MainModel
	LoginInput    textinput.Model
	PasswordInput textinput.Model
	ErrMessage    string
	step          int // шаги для последовательных действий (1 - первое, 2 - второе)
	stepMax       int // всего шагов в последовательности действий
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
		step:          stepInit,
		stepMax:       1,
	}
}

// ValidateScreenData проверяет и корректирует данные для текущего экрана.
func (s *RegisterScreen) ValidateScreenData() {
	s.step = stepInit
}

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

// Update описывает логику работы с командами для текущего окна.
func (s *RegisterScreen) Update(msg tea.Msg) (*MainModel, tea.Cmd) {
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

			ctx := context.Background()

			s.initAction(ctx, login, password)

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

func (s *RegisterScreen) initAction(inctx context.Context, login, password string) {
	ctx, cancelFn := context.WithCancel(inctx)

	s.step = 0
	s.stepMax = 2

	loadScreen := s.mainModel.screenLoading

	loadScreen.title = "Login user"
	loadScreen.desc = "Login user by email and password"
	loadScreen.percent = 0
	loadScreen.status = "Send request for login..."
	loadScreen.OnProgress = func(_ float64, _ string) tea.Cmd {
		return s.actionCmd(ctx, login, password)
	}
	loadScreen.OnDone = func(payload any) {
		s.mainModel.currentUser = &user{Login: login}

		nextScreen := s.mainModel.screenPassList

		items, _ := payload.([]service.Password)
		if items != nil {
			nextScreen.Items = items
		} else {
			nextScreen.Items = []service.Password{}
		}

		s.mainModel.SetCurrentScreen(nextScreen)
	}

	prevScreen := s.mainModel.screenLogin

	loadScreen.OnCancel = func() {
		cancelFn()

		prevScreen.ErrMessage = "operation canceled"

		s.mainModel.SetCurrentScreen(prevScreen)
	}
	loadScreen.OnError = func(err error) {
		prevScreen.ErrMessage = err.Error()

		s.mainModel.SetCurrentScreen(prevScreen)
	}

	s.mainModel.SetCurrentScreen(loadScreen)
}

func (s *RegisterScreen) actionCmd(
	ctx context.Context,
	login string,
	pass string,
) tea.Cmd {
	return func() tea.Msg {
		switch s.step {
		case stepInit:
			s.step = stepOne

			return LoadingProgressMsg{
				Percent: float64(s.step-1) / float64(s.stepMax),
				Status:  "Register user by email and password…",
			}

		case 1:
			if err := s.mainModel.service.Login(ctx, login, pass); err != nil {
				return LoadingDoneMsg{
					Err:     fmt.Errorf("authorization: %w", err),
					Payload: nil,
				}
			}

			s.step = stepTwo

			return LoadingProgressMsg{
				Percent: float64(s.step-1) / float64(s.stepMax),
				Status:  "Loading data…",
			}

		case stepTwo:
			items, err := s.mainModel.service.GetPasswords(ctx)
			if err != nil {
				return LoadingDoneMsg{
					Payload: nil,
					Err:     fmt.Errorf("loading data: %w", err),
				}
			}

			return LoadingDoneMsg{
				Payload: items,
				Err:     nil,
			}
		}

		return LoadingDoneMsg{
			Payload: nil,
			Err:     nil,
		}
	}
}
