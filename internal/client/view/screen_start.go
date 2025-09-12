package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type StartScreen struct {
	mainModel *model
	Index     int
	Items     []string
}

func NewStartScreen(mod *model) *StartScreen {
	return &StartScreen{
		mainModel: mod,
		Index:     0,
		Items: []string{
			"Login",
			"Register",
			"Quit",
		},
	}
}

func (s *StartScreen) LoadScreen(fnc func()) IScreen {
	s.Index = 0

	if fnc != nil {
		fnc()
	}

	min := 0
	if s.Index < min {
		s.Index = min
	}

	max := len(s.Items) - 1
	if s.Index > max {
		s.Index = max
	}

	return s
}

func (s *StartScreen) String() string {
	view := "\n[Menu] Select action:\n"

	for index := range s.Items {
		cursor := " "
		if index == s.Index {
			cursor = ">"
		}

		view += fmt.Sprintf("%s %s\n", cursor, s.Items[index])
	}

	return view
}

func (s *StartScreen) GetHints() []Hint {
	return []Hint{
		{"Select", []string{KeyEnter}},
		{"Switch", []string{KeyTab}},
		{"Next", []string{KeyDown}},
		{"Previous", []string{KeyUp}},
		{"Quit", []string{KeyEscape, KeyQuit}},
	}
}

func (s *StartScreen) Action(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case KeyEscape, KeyQuit:
		return s.mainModel, tea.Quit

	case KeyEnter:
		if s.Items[s.Index] == "Login" {
			s.mainModel.screenCurrent = s.mainModel.screenLogin.LoadScreen(func() {
				s.mainModel.screenLogin.ErrMessage = "WOW"
			})

			return s.mainModel, nil
		}

		if s.Items[s.Index] == "Register" {
			s.mainModel.screenCurrent = s.mainModel.screenRegister.LoadScreen(nil)

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

	case KeyUp:
		if s.Index > 0 {
			s.Index--
		}

		return s.mainModel, nil

	case KeyDown:
		if s.Index < len(s.Items)-1 {
			s.Index++
		}

		return s.mainModel, nil
	}

	return s.mainModel, nil
}
