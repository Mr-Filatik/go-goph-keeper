// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// RegisterScreen описывает экран рагистрации и необходимые ему данные.
type RegisterScreen struct {
	mainModel     *ViewModel
	LoginInput    textinput.Model
	PasswordInput textinput.Model
	ErrMessage    string
}

// NewRegisterScreen создаёт новый экзепляр *RegisterScreen.
func NewRegisterScreen(mod *ViewModel) *RegisterScreen {
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

func (s *RegisterScreen) LoadScreen(fnc func()) {
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
		{"Switch", []string{KeyTab, KeyDown, KeyUp}},
		{"Back", []string{KeyEscape}},
		{"Quit", []string{KeyQuit}},
	}
}

// Action описывает логику работы с командами для текущего окна.
func (s *RegisterScreen) Action(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, isKey := msg.(tea.KeyMsg)
	if isKey {
		switch key.String() {
		case KeyQuit:
			return s.mainModel, tea.Quit

		case KeyEscape:
			s.mainModel.screenStart.LoadScreen(nil)

			return s.mainModel, nil

		case KeyEnter:
			if s.LoginInput.Value() == "" || s.PasswordInput.Value() == "" {
				s.ErrMessage = "login and password are required"

				return s.mainModel, nil
			}

			s.mainModel.screenPassList.LoadScreen(func() {
				s.mainModel.currentUser = &user{
					Login: s.LoginInput.Value(),
				}
			})
			// s.ErrMessage = "happy (" + s.LoginInput.Value() + ") (" + s.PasswordInput.Value() + ")"

			return s.mainModel, nil

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
