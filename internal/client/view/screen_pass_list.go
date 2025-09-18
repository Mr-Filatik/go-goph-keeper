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
	step       int // шаги для последовательных действий (1 - первое, 2 - второе)
	stepMax    int // всего шагов в последовательности действий
}

// NewPasswordListScreen создаёт новый экзепляр *PasswordListScreen.
func NewPasswordListScreen(mod *MainModel) *PasswordListScreen {
	return &PasswordListScreen{
		mainModel:  mod,
		Index:      0,
		Items:      []service.Password{}, // + "[Add new password]"
		ErrMessage: "",
		step:       stepInit,
		stepMax:    1,
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

	s.step = stepInit
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
				}

				s.mainModel.SetCurrentScreen(screen)

				return s.mainModel, nil
			}

			ctx := context.Background()

			s.initAction(ctx, s.Items[s.Index].ID)

			return s.mainModel, s.actionCmd(ctx, s.Items[s.Index].ID)
		}
	}

	return s.mainModel, nil
}

func (s *PasswordListScreen) initAction(inctx context.Context, passID string) {
	ctx, cancelFn := context.WithCancel(inctx)

	s.step = 0
	s.stepMax = 1

	loadScreen := s.mainModel.screenLoading

	loadScreen.title = "Load secret data"
	loadScreen.desc = "Load secret data for " + s.Items[s.Index].Title
	loadScreen.percent = 0
	loadScreen.status = "Send request for load secret data..."
	loadScreen.OnProgress = func(_ float64, _ string) tea.Cmd {
		return s.actionCmd(ctx, passID)
	}
	loadScreen.OnDone = func(payload any) {
		nextScreen := s.mainModel.screenPassDetails

		pass, _ := payload.(string)
		s.Items[s.Index].Password = pass

		nextScreen.Item = &s.Items[s.Index]

		s.mainModel.SetCurrentScreen(nextScreen)
	}

	prevScreen := s.mainModel.screenPassList

	loadScreen.OnCancel = func() {
		cancelFn()

		prevScreen.ErrMessage = textOperationCanceled

		s.mainModel.SetCurrentScreen(prevScreen)
	}
	loadScreen.OnError = func(err error) {
		prevScreen.ErrMessage = err.Error()

		s.mainModel.SetCurrentScreen(prevScreen)
	}

	s.mainModel.SetCurrentScreen(loadScreen)
}

func (s *PasswordListScreen) actionCmd(
	ctx context.Context,
	passID string,
) tea.Cmd {
	return func() tea.Msg {
		switch s.step {
		case stepInit:
			s.step = stepOne

			return LoadingProgressMsg{
				Percent: float64(s.step-1) / float64(s.stepMax),
				Status:  "Load secret data…",
			}

		case 1:
			data, err := s.mainModel.service.GetPassword(ctx, passID)
			if err != nil {
				return LoadingDoneMsg{
					Err:     fmt.Errorf("get password info: %w", err),
					Payload: nil,
				}
			}

			return LoadingDoneMsg{
				Payload: data,
				Err:     nil,
			}
		}

		return LoadingDoneMsg{
			Payload: nil,
			Err:     nil,
		}
	}
}
