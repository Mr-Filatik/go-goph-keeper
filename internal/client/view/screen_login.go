package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type LoginScreen struct {
	mainModel *model
	Index     int
	Items     []string
}

func NewLoginScreen(mod *model) *LoginScreen {
	return &LoginScreen{
		mainModel: mod,
		Index:     0,
		Items: []string{
			"Login",
			"Register",
			"Quit",
		},
	}
}

func (s *LoginScreen) String() string {
	view := "\n[Login] Select action:\n"

	for index := range s.Items {
		cursor := " "
		if index == s.Index {
			cursor = ">"
		}

		view += fmt.Sprintf("%s %s\n", cursor, s.Items[index])
	}

	return view
}

func (s *LoginScreen) GetHints() []Hint {
	return []Hint{
		{"Select", []string{KeyEnter}},
		{"Switch", []string{KeyTab}},
		{"Next", []string{KeyDown, KeyDownWASD, KeyNext}},
		{"Previous", []string{KeyUp, KeyUpWASD, KeyPrev}},
		{"Quit", []string{KeyEscape, KeyEscapeShort, KeyQuit, KeyQuitShort}},
	}
}

func (s *LoginScreen) Action(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case KeyEscape, KeyEscapeShort, KeyQuit, KeyQuitShort:
		return s.mainModel, tea.Quit

	case KeyEnter:
		if s.Items[s.Index] == "Login" {
			s.mainModel.screenCurrent = s.mainModel.screenLogin

			return s.mainModel, nil
		}

		if s.Items[s.Index] == "Register" {
			s.mainModel.screenCurrent = s.mainModel.screenRegister

			return s.mainModel, nil
		}

		if s.Items[s.Index] == "Quit" {
			return s.mainModel, tea.Quit
		}

		return s.mainModel, nil

	case KeyTab:
		if s.Index < len(s.Items)-1 {
			s.Index++
		} else {
			s.Index = 0
		}

		return s.mainModel, nil

	case KeyUp, KeyUpWASD, KeyPrev:
		if s.Index > 0 {
			s.Index--
		}

		return s.mainModel, nil

	case KeyDown, KeyDownWASD, KeyNext:
		if s.Index < len(s.Items)-1 {
			s.Index++
		}

		return s.mainModel, nil
	}

	return s.mainModel, nil
}
