// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mr-filatik/go-goph-keeper/internal/client/service"
)

// MainModel описывает модель для отображения данных в консоли.
type MainModel struct {
	currentUser *user

	screenStart       *StartScreen
	screenLogin       *LoginScreen
	screenRegister    *RegisterScreen
	screenPassList    *PasswordListScreen
	screenPassDetails *PasswordDetailsScreen
	screenPassEdit    *PasswordEditScreen
	screenLoading     *LoadingScreen

	screenCurrent IScreen

	service service.IService

	// Глобальные данные, которые используются на всех окнах приложения.

	appName      string
	buildVersion string
	buildDate    string
	buildCommit  string
}

// NewMainModel создаёт новый экземпляр *MainModel.
func NewMainModel(serv service.IService) *MainModel {
	mod := &MainModel{
		appName:           "N/A",
		buildVersion:      "N/A",
		buildDate:         "N/A",
		buildCommit:       "N/A",
		currentUser:       nil,
		screenCurrent:     nil,
		screenStart:       nil,
		screenLogin:       nil,
		screenRegister:    nil,
		screenPassList:    nil,
		screenPassDetails: nil,
		screenPassEdit:    nil,
		screenLoading:     nil,
		service:           serv,
	}

	mod.screenStart = NewStartScreen(mod)
	mod.screenLogin = NewLoginScreen(mod)
	mod.screenRegister = NewRegisterScreen(mod)
	mod.screenPassList = NewPasswordListScreen(mod)
	mod.screenPassDetails = NewPasswordDetailsScreen(mod)
	mod.screenPassEdit = NewPasswordEditScreen(mod)
	mod.screenLoading = NewLoadingScreen(mod)

	mod.screenCurrent = mod.screenStart

	return mod
}

// Start запускает программу для отображения интерфейса.
func (m *MainModel) Start() error {
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return fmt.Errorf("run tea program: %w", err)
	}

	return nil
}

// SetCurrentScreen устанавливает новое текущее окно.
func (m *MainModel) SetCurrentScreen(screen IScreen) {
	if screen != nil {
		m.screenCurrent = screen
		m.screenCurrent.ValidateScreenData()
	}
}

// Init содержит действия, которые будут выполнены на этапе инициализации объекта.
func (m *MainModel) Init() tea.Cmd { return nil }

// Update реагирует на любые действия пользователя и изменения.
//
//nolint:ireturn // Bubble Tea требует возвращать tea.Model по интерфейсу
func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.screenCurrent != nil {
		return m.screenCurrent.Update(msg)
	}

	return m, nil
}

// View формирует вывод информации в консоль.
func (m *MainModel) View() string {
	view := m.viewHeader()

	if m.screenCurrent != nil {
		view += m.screenCurrent.String()
	}

	view += "\n" + m.viewFooter() + "\n"

	return view
}

// ExitToStartScreen отркрывает стартовый экран и удаляет авторизацию пользователя.
func (m *MainModel) ExitToStartScreen(ctx context.Context) (*MainModel, tea.Cmd) {
	m.SetCurrentScreen(m.screenStart)

	err := m.service.Logout(ctx)
	_ = err

	m.currentUser = nil

	return m, nil
}

func (m *MainModel) viewHeader() string {
	parts := []string{
		addLine(),
	}

	if m.currentUser != nil {
		parts = append(parts, m.appName+" [User] Login: "+m.currentUser.Login+".")
	} else {
		parts = append(parts, m.appName)
	}

	parts = append(parts, addLine(), "Hints:")

	if m.screenCurrent != nil {
		hints := m.screenCurrent.GetHints()
		if len(hints) != 0 {
			val := ""
			for _, hint := range hints {
				val += hint.actionName + " [" + strings.Join(hint.buttons, "/") + "] "
			}

			parts = append(parts, val, addLine()+"\n")
		}
	}

	return strings.Join(parts, "\n")
}

func (m *MainModel) viewFooter() string {
	return strings.Join([]string{
		addLine(),
		"[Build] Version: " + m.buildVersion + ", Date: " + m.buildDate + ".",
		addLine() + "\n",
	}, "\n")
}

const (
	lineWidth  = 60
	lineSymbol = "─"
)

func addLine() string {
	b := make([]byte, 0, len(lineSymbol)*lineWidth)
	for range lineWidth {
		b = append(b, lineSymbol...)
	}

	return string(b)
}
