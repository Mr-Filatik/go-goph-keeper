// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// PasswordListScreen описывает экран всех паролей и необходимые ему данные.
type PasswordListScreen struct {
	mainModel *ViewModel
	Index     int
	Items     []*item
}

// NewPasswordListScreen создаёт новый экзепляр *PasswordListScreen.
func NewPasswordListScreen(mod *ViewModel) *PasswordListScreen {
	return &PasswordListScreen{
		mainModel: mod,
		Index:     0,
		Items: []*item{ // "[Add new password]"
			{"GitHub", "Personal account", "vlad", "ghp_example_password", "2FA: TOTP in Authy"},
			{"GMail", "Work", "vladislav", "gmail_app_password", "App password only"},
			{"AWS", "Prod account", "admin", "supersecretkey", "Use IAM roles"},
		},
	}
}

func (s *PasswordListScreen) LoadScreen(fnc func()) {
	s.mainModel.screenCurrent = s

	s.Index = 0
	// s.Items = []*item{
	// 	{"GitHub", "Personal account", "vlad", "ghp_example_password", "2FA: TOTP in Authy"},
	// 	{"GMail", "Work", "vladislav", "gmail_app_password", "App password only"},
	// 	{"AWS", "Prod account", "admin", "supersecretkey", "Use IAM roles"},
	// }

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
func (s *PasswordListScreen) String() string {
	view := "\n[Password List] Select password:\n"

	for index := range s.Items {
		cursor := " "
		if index == s.Index {
			cursor = ">"
		}

		view += fmt.Sprintf("%s %s\n", cursor, s.Items[index].Title)
	}

	return view
}

// GetHints выводит подсказки по управлению для текущего окна.
func (s *PasswordListScreen) GetHints() []Hint {
	return []Hint{
		{"Login", []string{KeyEnter}},
		{"Next", []string{KeyDown}},
		{"Previous", []string{KeyUp}},
		{"Back", []string{KeyEscape}},
		{"Quit", []string{KeyQuit}},
	}
}

// Action описывает логику работы с командами для текущего окна.
func (s *PasswordListScreen) Action(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, isKey := msg.(tea.KeyMsg); isKey {
		switch key.String() {
		case KeyQuit:
			return s.mainModel, tea.Quit

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
			s.mainModel.screenPassDetails.LoadScreen(func() {
				s.mainModel.screenPassDetails.Item = s.Items[s.Index]
			})

			return s.mainModel, nil

			// s.mainModel.selected = &s.mainModel.items[s.mainModel.listIndex]
			// s.mainModel.currentScreen = screenPasswordDetails
			// s.mainModel.statusMsg = "Opened details"
		}
	}

	return s.mainModel, nil
}
