// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mr-filatik/go-goph-keeper/internal/client/service"
)

// PasswordEditScreen описывает экран редактирования пароля и необходимые ему данные.
type PasswordEditScreen struct {
	mainModel   *MainModel
	Index       int
	Items       []string
	Item        *service.Password
	InfoMessage string
	IsCreate    bool
}

// NewPasswordEditScreen создаёт новый экзепляр *PasswordEditScreen.
func NewPasswordEditScreen(mod *MainModel) *PasswordEditScreen {
	return &PasswordEditScreen{
		mainModel: mod,
		Index:     0,
		Items: []string{
			"Back to details",
			"Edit", // если создаётся заново, то кнопки другие Add вместо Edit Remove
			"Remove",
		},
		Item:        nil,
		InfoMessage: "",
		IsCreate:    false,
	}
}

// ValidateScreenData проверяет и корректирует данные для текущего экрана.
func (s *PasswordEditScreen) ValidateScreenData() {
	minLimit := 0
	if s.Index < minLimit {
		s.Index = minLimit
	}

	maxLimit := len(s.Items) - 1
	if s.Index > maxLimit {
		s.Index = maxLimit
	}

	if s.IsCreate {
		s.Items = []string{
			"Create",
			"Back to details",
		}
	} else {
		s.Items = []string{
			"Back to details",
			"Edit",
			"Remove",
		}
	}
}

// String выводит окно и его содержимое в виде строки.
func (s *PasswordEditScreen) String() string {
	view := "\n[Password Edit] Edit:\n"

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
func (s *PasswordEditScreen) GetHints() []Hint {
	return []Hint{
		{"Select", []string{KeyEnter}},
		{"Switch", []string{KeyTab}},
		{"Next", []string{KeyDown}},
		{"Previous", []string{KeyUp}},
		{"Back", []string{KeyEscape}},
		// ctrl+c = save
	}
}

// Update описывает логику работы с командами для текущего окна.
func (s *PasswordEditScreen) Update(msg tea.Msg) (*MainModel, tea.Cmd) {
	if key, isKey := msg.(tea.KeyMsg); isKey {
		switch key.String() {
		case KeyEscape:
			s.actionBackToDetails()

			return s.mainModel, nil

		case KeyUp:
			s.Index = indexPrev(s.Index)

			return s.mainModel, nil

		case KeyDown:
			s.Index = indexNext(s.Index, len(s.Items))

			return s.mainModel, nil

		case KeyEnter:
			s.enter()

			return s.mainModel, nil

		case KeyCopy:
			return s.mainModel, nil
		}
	}

	return s.mainModel, nil
}

func (s *PasswordEditScreen) enter() {
	if s.Items[s.Index] == "Create" {
		// logic
		_ = s.IsCreate
	}

	if s.Items[s.Index] == "Back to details" {
		if s.IsCreate {
			screen := s.mainModel.screenPassList

			s.mainModel.SetCurrentScreen(screen)
		} else {
			s.actionBackToDetails()
		}
	}

	if s.Items[s.Index] == "Edit" {
		// logic
		_ = s.IsCreate
	}

	if s.Items[s.Index] == "Remove" {
		// logic
		_ = s.IsCreate
	}
}

func (s *PasswordEditScreen) actionBackToDetails() {
	screen := s.mainModel.screenPassDetails

	screen.Item = s.Item

	s.mainModel.SetCurrentScreen(screen)
}
