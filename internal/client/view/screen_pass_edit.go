// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
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
	ErrMessage  string
	IsCreate    bool
	step        int // шаги для последовательных действий (1 - первое, 2 - второе)
	stepMax     int // всего шагов в последовательности действий

	editMode  bool            // false: выбор поля/кнопок; true: редактирование поля
	editField string          // какое поле редактируем
	input     textinput.Model // textinput для редактирования
}

// NewPasswordEditScreen создаёт новый экзепляр *PasswordEditScreen.
func NewPasswordEditScreen(mod *MainModel) *PasswordEditScreen {
	input := textinput.New()
	input.CharLimit = 256

	screen := &PasswordEditScreen{
		mainModel: mod,
		Index:     0,
		Items: []string{
			"Back to details",
			"Edit", // если создаётся заново, то кнопки другие Add вместо Edit Remove
			"Remove",
		},
		Item:        nil,
		InfoMessage: "",
		ErrMessage:  "",
		IsCreate:    false,
		input:       input,
		step:        stepInit,
		stepMax:     1,
		editMode:    false,
		editField:   "",
	}
	screen.rebuildMenu()

	return screen
}

// ValidateScreenData проверяет и корректирует данные для текущего экрана.
func (s *PasswordEditScreen) ValidateScreenData() {
	s.rebuildMenu()
}

// String выводит окно и его содержимое в виде строки.
//
//nolint:cyclop
func (s *PasswordEditScreen) String() string {
	view := "\n[Password Edit] "

	if s.IsCreate {
		view += "Create new item:\n"
	} else {
		view += "Edit item:\n"
	}

	if s.Item == nil {
		view += "\nNo data...\n"
	} else {
		item := s.Item

		view += fmt.Sprintf("Type: %s\n", safeType(item.Type))
		view += fmt.Sprintf("Title: %s\n", item.Title)
		view += fmt.Sprintf("Description: %s\n", item.Description)

		// Поля Login/Password показываем для типа login
		if item.Type == service.PasswordTypeLogin {
			view += fmt.Sprintf("Login: %s\n", item.Login)
			view += "Password: [hidden]\n"
		} else {
			view += "(no login/password for this type)\n"
		}
	}

	view += "\n[Edit] Select field/action:\n"

	for index := range s.Items {
		cursor := " "
		if index == s.Index {
			cursor = ">"
		}

		line := s.Items[index]
		if s.editMode && line == s.editField {
			line += " [editing]"
		}

		view += fmt.Sprintf("%s %s\n", cursor, line)
	}

	if s.editMode {
		view += "\n[Input]: " + s.input.View() + "\n"
	}

	if s.InfoMessage != "" {
		view += "\n[INFO]: " + s.InfoMessage + "\n"
	}

	if s.ErrMessage != "" {
		view += "\n[ERROR]: " + s.ErrMessage + "\n"
	}

	return view
}

// GetHints выводит подсказки по управлению для текущего окна.
func (s *PasswordEditScreen) GetHints() []Hint {
	if s.editMode {
		return []Hint{
			{"Apply", []string{KeyEnter}},
			{"Cancel", []string{KeyEscape}},
		}
	}

	return []Hint{
		{"Select", []string{KeyEnter}},
		{"Switch", []string{KeyTab, KeyDown, KeyUp}},
		{"Back", []string{KeyEscape}},
	}
}

// Update описывает логику работы с командами для текущего окна.
//
//nolint:cyclop
func (s *PasswordEditScreen) Update(msg tea.Msg) (*MainModel, tea.Cmd) {
	// Режим редактирования: отдаём события в textinput
	if s.editMode {
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case KeyEscape:
				// Выходим из редактирования без сохранения
				s.editMode = false
				s.InfoMessage = "edit canceled"

				return s.mainModel, nil
			case KeyEnter:
				// Применим введённое значение
				s.applyEdit()
				s.editMode = false

				return s.mainModel, nil
			}
		}

		var cmd tea.Cmd
		s.input, cmd = s.input.Update(msg)

		return s.mainModel, cmd
	}

	// Обычный режим: выбор пунктов
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case KeyEscape:
			if s.IsCreate {
				// Назад в список
				screen := s.mainModel.screenPassList
				s.mainModel.SetCurrentScreen(screen)

				return s.mainModel, nil
			}

			s.actionBackToDetails()

			return s.mainModel, nil

		case KeyUp:
			s.Index = indexPrev(s.Index)

			return s.mainModel, nil

		case KeyDown:
			s.Index = indexNext(s.Index, len(s.Items))

			return s.mainModel, nil

		case KeyTab:
			s.Index = indexSwitch(s.Index, len(s.Items))

			return s.mainModel, nil

		case KeyEnter:
			return s.mainModel, s.enter()
		}
	}

	return s.mainModel, nil
}

// rebuildMenu перестраивает список пунктов в зависимости от IsCreate и типа записи.
func (s *PasswordEditScreen) rebuildMenu() {
	if s.IsCreate {
		s.Items = []string{
			"Title",
			"Description",
			"Type",
			"Login",
			"Password",
			"Create",
			"Back to list",
		}
	} else {
		s.Items = []string{
			"Title",
			"Description",
			"Type",
			"Login",
			"Password",
			"Save",
			"Remove",
			"Back to details",
		}
	}

	if s.Index < 0 {
		s.Index = 0
	}

	if s.Index >= len(s.Items) {
		s.Index = len(s.Items) - 1
	}
}

// enter — обработка выбранного пункта меню (в обычном режиме).
//
//nolint:cyclop
func (s *PasswordEditScreen) enter() tea.Cmd {
	if s.Item == nil {
		return nil
	}

	choice := s.Items[s.Index]
	switch choice {
	//nolint:goconst
	case "Title", "Description", "Login", "Password":
		// Для типа != login логин/пароль пропустим
		if (choice == "Login" || choice == "Password") && s.Item.Type != service.PasswordTypeLogin {
			s.InfoMessage = "field not available for current type"

			return nil
		}

		s.startEdit(choice)

		return nil

	case "Type":
		s.cycleType()

		return nil

	case "Create", "Save":
		if s.IsCreate {
			return s.createItem()
		}

		return s.saveItem()

	case "Remove":
		if s.IsCreate {
			return nil
		}

		// В прототипе не нужно опустим
		s.InfoMessage = "remove not implemented"

		return nil

	case "Back to list":
		screen := s.mainModel.screenPassList
		s.mainModel.SetCurrentScreen(screen)

		return nil

	case "Back to details":
		if s.IsCreate {
			screen := s.mainModel.screenPassList
			s.mainModel.SetCurrentScreen(screen)
		} else {
			s.actionBackToDetails()
		}

		return nil
	}

	return nil
}

// startEdit включает режим редактирования для поля.
func (s *PasswordEditScreen) startEdit(field string) {
	s.editMode = true
	s.editField = field

	// начальное значение в input
	switch field {
	case "Title":
		s.input.SetValue(s.Item.Title)
	case "Description":
		s.input.SetValue(s.Item.Description)
	case "Login":
		s.input.SetValue(s.Item.Login)
	case "Password":
		s.input.SetValue(s.Item.Password)
		s.input.EchoMode = textinput.EchoPassword
		s.input.EchoCharacter = '•'
	}

	s.input.Focus()
}

// applyEdit применяет введённый текст к выбранному полю.
func (s *PasswordEditScreen) applyEdit() {
	val := s.input.Value()

	switch s.editField {
	case "Title":
		s.Item.Title = val
	case "Description":
		s.Item.Description = val
	case "Login":
		s.Item.Login = val
	case "Password":
		s.Item.Password = val
	}

	s.InfoMessage = s.editField + " updated"
	// вернуть обычный режим ввода
	s.input.Blur()
	s.input.EchoMode = textinput.EchoNormal
	s.input.SetValue("")
}

// cycleType циклически меняет тип записи.
func (s *PasswordEditScreen) cycleType() {
	next := map[service.ItemType]service.ItemType{
		service.PasswordTypeLogin:  service.PasswordTypeText,
		service.PasswordTypeText:   service.PasswordTypeBinary,
		service.PasswordTypeBinary: service.PasswordTypeCard,
		service.PasswordTypeCard:   service.PasswordTypeLogin,
	}
	s.Item.Type = next[s.Item.Type]
	s.InfoMessage = "type: " + string(s.Item.Type)
}

// createItem вызывает AddPassword при IsCreate=true.
func (s *PasswordEditScreen) createItem() tea.Cmd {
	item := *s.Item // копия

	return func() tea.Msg {
		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()

		newID, err := s.mainModel.service.AddPassword(ctx, item)
		if err != nil {
			s.InfoMessage = "create error: " + err.Error()

			return nil
		}

		s.InfoMessage = "created with ID: " + newID
		// после создания — в список
		screen := s.mainModel.screenPassList
		s.mainModel.SetCurrentScreen(screen)

		return nil
	}
}

// saveItem вызывает ChangePassword при редактировании существующей записи.
func (s *PasswordEditScreen) saveItem() tea.Cmd {
	item := *s.Item // копия

	return func() tea.Msg {
		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()

		if err := s.mainModel.service.ChangePassword(ctx, item); err != nil {
			s.InfoMessage = "save error: " + err.Error()

			return nil
		}

		s.InfoMessage = "saved"
		// назад в детали
		s.actionBackToDetails()

		return nil
	}
}

func (s *PasswordEditScreen) actionBackToDetails() {
	screen := s.mainModel.screenPassDetails
	screen.Item = s.Item
	s.mainModel.SetCurrentScreen(screen)
}

func safeType(t service.ItemType) service.ItemType {
	if t == "" {
		return service.PasswordTypeLogin
	}

	return t
}
