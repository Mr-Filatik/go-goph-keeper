// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// StartScreen описывает начальный экран приложения и необходимые ему данные.
type StartScreen struct {
	mainModel *ViewModel
	Index     int
	Items     []string
}

// NewStartScreen создаёт новый экзепляр *StartScreen.
func NewStartScreen(mod *ViewModel) *StartScreen {
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

func (s *StartScreen) LoadScreen(fnc func()) {
	s.mainModel.screenCurrent = s

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
}

// String выводит окно и его содержимое в виде строки.
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

// GetHints выводит подсказки по управлению для текущего окна.
func (s *StartScreen) GetHints() []Hint {
	return []Hint{
		{"Select", []string{KeyEnter}},
		{"Switch", []string{KeyTab}},
		{"Next", []string{KeyDown}},
		{"Previous", []string{KeyUp}},
		{"Quit", []string{KeyEscape, KeyQuit}},
	}
}

// Action описывает логику работы с командами для текущего окна.
func (s *StartScreen) Action(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, isKey := msg.(tea.KeyMsg); isKey {
		switch key.String() {
		case KeyEscape, KeyQuit:
			return s.mainModel, tea.Quit

		case KeyEnter:
			if s.Items[s.Index] == "Login" {
				s.mainModel.screenLogin.LoadScreen(nil)

				return s.mainModel, nil
			}

			if s.Items[s.Index] == "Register" {
				s.mainModel.screenRegister.LoadScreen(nil)

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
	}

	return s.mainModel, nil
}
