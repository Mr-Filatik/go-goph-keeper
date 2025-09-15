// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ViewModel описывает модель для отображения данных в консоли.
type ViewModel struct {
	currentUser *user

	screenStart       *StartScreen
	screenLogin       *LoginScreen
	screenRegister    *RegisterScreen
	screenPassList    *PasswordListScreen
	screenPassDetails *PasswordDetailsScreen
	screenPassEdit    *PasswordEditScreen
	screenLoading     *LoadingScreen

	screenCurrent IScreen

	// Глобальные данные, которые используются на всех окнах приложения.

	appName      string
	buildVersion string
	buildDate    string
	buildCommit  string
}

func newViewModel() *ViewModel {
	mod := &ViewModel{
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

// SetCurrentScreen устанавливает новое текущее окно.
func (m *ViewModel) SetCurrentScreen(screen IScreen) {
	if screen != nil {
		m.screenCurrent = screen
	}
}

func (m *ViewModel) Init() tea.Cmd { return nil }

// Кнопки для управления UI.
const (
	// Элементы управления.
	KeyTab  = "tab"
	KeyUp   = "up"
	KeyDown = "down"

	// Элементы действий.
	KeyEnter = "enter"
	KeyCopy  = "ctrl+c"

	// Элементы для выхода.
	KeyEscape = "esc"
	KeyQuit   = "ctrl+q"
)

func (m *ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.screenCurrent != nil {
		return m.screenCurrent.Action(msg)
	}

	return m, nil
}

func (m *ViewModel) View() string {
	view := m.header()

	if m.screenCurrent != nil {
		view += m.screenCurrent.String()
	}

	view += "\n" + m.footer() + "\n"

	return view
}

const (
	lineWidth  = 60
	lineSymbol = "─"
)

func (m *ViewModel) header() string {
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

func (m *ViewModel) footer() string {
	return strings.Join([]string{
		addLine(),
		"[Build] Version: " + m.buildVersion + ", Date: " + m.buildDate + ".",
		addLine() + "\n",
	}, "\n")

	// if m.statusMsg == "" {
	// 	return stringsRepeat("─", 60)
	// }
	// return fmt.Sprintf("%s\n%s", m.statusMsg, stringsRepeat("─", 60))
}

func addLine() string {
	b := make([]byte, 0, len(lineSymbol)*lineWidth)
	for range lineWidth {
		b = append(b, lineSymbol...)
	}
	return string(b)
}

func Start() {
	m := newViewModel()
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
