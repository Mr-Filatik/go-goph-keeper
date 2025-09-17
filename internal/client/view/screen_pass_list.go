// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mr-filatik/go-goph-keeper/internal/client/service"
)

// PasswordListScreen описывает экран всех паролей и необходимые ему данные.
type PasswordListScreen struct {
	mainModel  *MainModel
	Index      int
	Items      []service.Password
	ErrMessage string
}

// NewPasswordListScreen создаёт новый экзепляр *PasswordListScreen.
func NewPasswordListScreen(mod *MainModel) *PasswordListScreen {
	return &PasswordListScreen{
		mainModel:  mod,
		Index:      0,
		Items:      []service.Password{}, // + "[Add new password]"
		ErrMessage: "",
	}
}

// ValidateScreenData проверяет и корректирует данные для текущего экрана.
func (s *PasswordListScreen) ValidateScreenData() {
	if len(s.Items) == 0 {
		s.Index = -1
	}

	if s.Index > len(s.Items)-1 {
		s.Index = len(s.Items) - 1
	}
}

// String выводит окно и его содержимое в виде строки.
func (s *PasswordListScreen) String() string {
	view := "\n[Password List] Select password:\n"

	for index := -1; index < len(s.Items); index++ {
		cursor := " "

		if index == s.Index {
			cursor = ">"
		}

		if index == -1 {
			view += fmt.Sprintf("%s %s\n", cursor, "Add new password")

			continue
		}

		view += fmt.Sprintf("%s %s\n", cursor, s.Items[index].Title)
	}

	if s.ErrMessage != "" {
		view += "\n[ERROR]: " + s.ErrMessage + "\n"
	}

	return view
}

// GetHints выводит подсказки по управлению для текущего окна.
func (s *PasswordListScreen) GetHints() []Hint {
	return []Hint{
		{"Login", []string{KeyEnter}},
		{"Next", []string{KeyDown}},
		{"Previous", []string{KeyUp}},
		{"Exit", []string{KeyEscape}},
	}
}

// Update описывает логику работы с командами для текущего окна.
func (s *PasswordListScreen) Update(msg tea.Msg) (*MainModel, tea.Cmd) {
	if key, isKey := msg.(tea.KeyMsg); isKey {
		switch key.String() {
		case KeyEscape:
			return s.mainModel.ExitToStartScreen(context.Background())

		case KeyUp:
			// дополнительный пункт "Add new password" с индексом -1.
			s.Index = indexPrevWithCustomLimit(s.Index, -1)

			return s.mainModel, nil

		case KeyDown:
			s.Index = indexNext(s.Index, len(s.Items))

			return s.mainModel, nil

		case KeyEnter:
			// дополнительный пункт "Add new password" с индексом -1.
			if s.Index == -1 {
				screen := s.mainModel.screenPassEdit

				screen.IsCreate = true
				screen.Item = &service.Password{
					ID:          "",
					Title:       "",
					Description: "",
					Login:       "",
					Password:    "",
					Notes:       "",
				}

				s.mainModel.SetCurrentScreen(screen)

				return s.mainModel, nil
			}

			screen := s.mainModel.screenPassDetails

			screen.Item = &s.Items[s.Index]

			s.mainModel.SetCurrentScreen(screen)

			return s.mainModel, nil
		}
	}

	return s.mainModel, nil
}
