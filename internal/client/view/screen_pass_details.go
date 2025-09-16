// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"fmt"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mr-filatik/go-goph-keeper/internal/client/service"
)

// PasswordDetailsScreen описывает экран детальной информации о пароле и необходимые ему данные.
type PasswordDetailsScreen struct {
	mainModel   *MainModel
	Index       int
	Items       []string
	Item        *service.Password
	InfoMessage string
}

// NewPasswordDetailsScreen создаёт новый экзепляр *PasswordDetailsScreen.
func NewPasswordDetailsScreen(mod *MainModel) *PasswordDetailsScreen {
	return &PasswordDetailsScreen{
		mainModel: mod,
		Index:     0,
		Items: []string{
			"Copy login",
			"Copy password",
			"Edit",
			"Back to list",
		},
		Item:        nil,
		InfoMessage: "",
	}
}

// ValidateScreenData проверяет и корректирует данные для текущего экрана.
func (s *PasswordDetailsScreen) ValidateScreenData() {
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
func (s *PasswordDetailsScreen) String() string {
	view := "\n[Password Details] Data:\n"

	if s.Item == nil {
		view += "\nNo data...\n"
	} else {
		item := s.Item
		view += fmt.Sprintf("Name: %s\nDecsription: %s\n", item.Title, item.Description)
		view += fmt.Sprintf("Login: %s\n", item.Login)
		view += "Password: [hidden] (press 'ctrl+c' to copy)\n"
	}

	view += "\n[Password Details] Select action:\n"

	for index := range s.Items {
		cursor := " "
		if index == s.Index {
			cursor = ">"
		}

		view += fmt.Sprintf("%s %s\n", cursor, s.Items[index])
	}

	if s.InfoMessage != "" {
		view += "\n[INFO]: " + s.InfoMessage + "\n"
	}

	return view
}

// GetHints выводит подсказки по управлению для текущего окна.
func (s *PasswordDetailsScreen) GetHints() []Hint {
	return []Hint{
		{"Select", []string{KeyEnter}},
		{"Switch", []string{KeyTab}},
		{"Next", []string{KeyDown}},
		{"Previous", []string{KeyUp}},
		{"Back", []string{KeyEscape}},
	}
}

// Action описывает логику работы с командами для текущего окна.
func (s *PasswordDetailsScreen) Action(msg tea.Msg) (*MainModel, tea.Cmd) {
	if key, isKey := msg.(tea.KeyMsg); isKey {
		switch key.String() {
		case KeyEscape:
			return s.actionBackToList()

		case KeyUp:
			s.Index = indexPrev(s.Index)

			return s.mainModel, nil

		case KeyDown:
			s.Index = indexNext(s.Index, len(s.Items))

			return s.mainModel, nil

		case KeyEnter:
			if s.Items[s.Index] == "Copy login" {
				if s.Item != nil {
					_ = clipboard.WriteAll(s.Item.Login)
					s.InfoMessage = "login copied to clipboard"
				}

				return s.mainModel, nil
			}

			if s.Items[s.Index] == "Copy password" {
				if s.Item != nil {
					_ = clipboard.WriteAll(s.Item.Password)
					s.InfoMessage = "password copied to clipboard"
				}

				return s.mainModel, nil
			}

			if s.Items[s.Index] == "Edit" {
				// s.mainModel.screenCurrent = s.mainModel.screenRegister.LoadScreen(nil)

				return s.mainModel, nil
			}

			if s.Items[s.Index] == "Back to list" {
				return s.actionBackToList()
			}

			return s.mainModel, nil

		case KeyCopy:
			if s.Item != nil {
				_ = clipboard.WriteAll(s.Item.Password)
				s.InfoMessage = "password copied to clipboard"
			}

			return s.mainModel, nil
		}
	}

	return s.mainModel, nil
}

func (s *PasswordDetailsScreen) actionBackToList() (*MainModel, tea.Cmd) {
	screen := s.mainModel.screenPassList

	s.mainModel.SetCurrentScreen(screen)

	// s.mainModel.screenPassList.LoadScreen(func() {
	// 	for i := 0; i < len(s.mainModel.screenPassList.Items); i++ {
	// 		if s.Item == s.mainModel.screenPassList.Items[i] {
	// 			s.mainModel.screenPassList.Index = i
	// 		}
	// 	}
	// })

	return s.mainModel, nil
}
