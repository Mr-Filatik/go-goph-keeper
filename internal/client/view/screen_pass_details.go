package view

import (
	"fmt"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

type PasswordDetailsScreen struct {
	mainModel   *model
	Index       int
	Items       []string
	Item        *item
	InfoMessage string
}

func NewPasswordDetailsScreen(mod *model) *PasswordDetailsScreen {
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

func (s *PasswordDetailsScreen) LoadScreen(fnc func()) IScreen {
	s.Index = 0
	s.Item = nil
	s.InfoMessage = ""

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

func (s *PasswordDetailsScreen) String() string {
	view := "\n[Password Details] Data:\n"

	if s.Item == nil {
		view += "\nNo data...\n"
	} else {
		item := s.Item
		view += fmt.Sprintf("Name: %s\nDecsription: %s\n", item.Title, item.Description)
		view += fmt.Sprintf("Login: %s\n", item.Username)
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

func (s *PasswordDetailsScreen) GetHints() []Hint {
	return []Hint{
		{"Select", []string{KeyEnter}},
		{"Switch", []string{KeyTab}},
		{"Next", []string{KeyDown}},
		{"Previous", []string{KeyUp}},
		{"Quit", []string{KeyEscape, KeyQuit}},
	}
}

func (s *PasswordDetailsScreen) Action(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case KeyQuit:
		return s.mainModel, tea.Quit

	case KeyEscape:
		return s.actionBackToList()

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

	case KeyEnter:
		if s.Items[s.Index] == "Copy login" {
			if s.Item != nil {
				_ = clipboard.WriteAll(s.Item.Username)
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

	return s.mainModel, nil
}

func (s *PasswordDetailsScreen) actionBackToList() (tea.Model, tea.Cmd) {
	s.mainModel.screenCurrent = s.mainModel.screenPassList.LoadScreen(func() {
		for i := 0; i < len(s.mainModel.screenPassList.Items); i++ {
			if s.Item == s.mainModel.screenPassList.Items[i] {
				s.mainModel.screenPassList.Index = i
			}
		}
	})

	return s.mainModel, nil
}
