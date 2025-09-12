package view

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginScreen struct {
	mainModel     *model
	LoginInput    textinput.Model
	PasswordInput textinput.Model
	ErrMessage    string
}

func NewLoginScreen(mod *model) *LoginScreen {
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

func (s *LoginScreen) LoadScreen(fnc func()) IScreen {
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

	return s
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

func (s *LoginScreen) Action(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case KeyQuit:
		return s.mainModel, tea.Quit

	case KeyEscape:
		s.mainModel.screenCurrent = s.mainModel.screenStart.LoadScreen(nil)

		return s.mainModel, nil

	case KeyEnter:
		if s.LoginInput.Value() == "" || s.PasswordInput.Value() == "" {
			s.ErrMessage = "login and password are required"

			return s.mainModel, nil
		}

		s.mainModel.screenCurrent = s.mainModel.screenPassList.LoadScreen(func() {
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

	var cmd tea.Cmd
	if s.LoginInput.Focused() {
		s.LoginInput, cmd = s.LoginInput.Update(key)
	} else {
		s.PasswordInput, cmd = s.PasswordInput.Update(key)
	}

	return s.mainModel, cmd
}
