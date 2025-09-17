// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// StartScreen описывает начальный экран приложения и необходимые ему данные.
type StartScreen struct {
	mainModel *MainModel
	Index     int
	Items     []string
}

// NewStartScreen создаёт новый экзепляр *StartScreen.
func NewStartScreen(mod *MainModel) *StartScreen {
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

// ValidateScreenData проверяет и корректирует данные для текущего экрана.
func (s *StartScreen) ValidateScreenData() {
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
		{"Quit", []string{KeyEscape}},
	}
}

// Update описывает логику работы с командами для текущего окна.
func (s *StartScreen) Update(msg tea.Msg) (*MainModel, tea.Cmd) {
	if key, isKey := msg.(tea.KeyMsg); isKey {
		switch key.String() {
		case KeyEscape:
			return s.mainModel, tea.Quit

		case KeyEnter:
			if s.Items[s.Index] == "Login" {
				s.mainModel.SetCurrentScreen(s.mainModel.screenLogin)

				return s.mainModel, nil
			}

			if s.Items[s.Index] == "Register" {
				s.mainModel.SetCurrentScreen(s.mainModel.screenRegister)

				return s.mainModel, nil
			}

			if s.Items[s.Index] == "Quit" {
				return s.mainModel, tea.Quit
			}

			return s.mainModel, nil

		case KeyTab:
			s.Index = indexSwitch(s.Index, len(s.Items))

			return s.mainModel, nil

		case KeyUp:
			s.Index = indexPrev(s.Index)

			return s.mainModel, nil

		case KeyDown:
			s.Index = indexNext(s.Index, len(s.Items))

			return s.mainModel, nil
		}
	}

	return s.mainModel, nil
}
