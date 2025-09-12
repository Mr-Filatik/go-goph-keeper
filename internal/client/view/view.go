package view

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type availableAction struct {
	name    string
	buttons []string
}

type model struct {
	currentUser *user

	// global/app info always shown in header
	appName      string
	buildVersion string
	buildDate    string
	buildCommit  string

	// loading screen data
	loadingInfo    string
	loadingPercent float64

	screenCurrent     IScreen
	screenStart       *StartScreen
	screenLogin       *LoginScreen
	screenRegister    *RegisterScreen
	screenPassList    *PasswordListScreen
	screenPassDetails *PasswordDetailsScreen
}

func initialModel() model {
	mod := &model{
		appName:      "AppName",
		buildVersion: "Version",
		buildDate:    "Build",
		buildCommit:  "Commit",
		currentUser:  nil,
	}

	mod.screenStart = NewStartScreen(mod)
	mod.screenLogin = NewLoginScreen(mod)
	mod.screenRegister = NewRegisterScreen(mod)
	mod.screenPassList = NewPasswordListScreen(mod)
	mod.screenPassDetails = NewPasswordDetailsScreen(mod)

	mod.screenCurrent = mod.screenStart

	return *mod
}

func (m model) Init() tea.Cmd { return nil }

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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, isKey := msg.(tea.KeyMsg); isKey {
		if m.screenCurrent != nil {
			return m.screenCurrent.Action(key)
		}
	}

	return m, nil
}

func (m model) View() string {
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

func (m model) header() string {
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

func (m model) footer() string {
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
	m := initialModel()
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

// type Model struct {
// 	choices  []string         // items on the to-do list
// 	cursor   int              // which to-do list item our cursor is pointing at
// 	selected map[int]struct{} // which to-do items are selected
// }

// func InitialModel() Model {
// 	return Model{
// 		// Our to-do list is a grocery list
// 		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

// 		// A map which indicates which choices are selected. We're using
// 		// the  map like a mathematical set. The keys refer to the indexes
// 		// of the `choices` slice, above.
// 		selected: make(map[int]struct{}),
// 	}
// }

// func Start() {
// 	program := tea.NewProgram(InitialModel())
// 	_, err := program.Run()
// 	if err != nil {
// 		fmt.Print("TEA ERROR")
// 	}
// }

// func (m Model) Init() tea.Cmd {
// 	// Just return `nil`, which means "no I/O right now, please."
// 	return nil
// }

// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {

// 	// Is it a key press?
// 	case tea.KeyMsg:

// 		// Cool, what was the actual key pressed?
// 		switch msg.String() {

// 		// These keys should exit the program.
// 		case "ctrl+c", "q":
// 			return m, tea.Quit

// 		// The "up" and "k" keys move the cursor up
// 		case "up", "k":
// 			if m.cursor > 0 {
// 				m.cursor--
// 			}

// 		// The "down" and "j" keys move the cursor down
// 		case "down", "j":
// 			if m.cursor < len(m.choices)-1 {
// 				m.cursor++
// 			}

// 		// The "enter" key and the spacebar (a literal space) toggle
// 		// the selected state for the item that the cursor is pointing at.
// 		case "enter", " ":
// 			_, ok := m.selected[m.cursor]
// 			if ok {
// 				delete(m.selected, m.cursor)
// 			} else {
// 				m.selected[m.cursor] = struct{}{}
// 			}
// 		}
// 	}

// 	// Return the updated model to the Bubble Tea runtime for processing.
// 	// Note that we're not returning a command.
// 	return m, nil
// }

// func (m Model) View() string {
// 	// The header
// 	s := "What should we buy at the market?\n\n"

// 	// Iterate over our choices
// 	for i, choice := range m.choices {

// 		// Is the cursor pointing at this choice?
// 		cursor := " " // no cursor
// 		if m.cursor == i {
// 			cursor = ">" // cursor!
// 		}

// 		// Is this choice selected?
// 		checked := " " // not selected
// 		if _, ok := m.selected[i]; ok {
// 			checked = "x" // selected!
// 		}

// 		// Render the row
// 		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
// 	}

// 	// The footer
// 	s += "\nPress q to quit.\n"

// 	// Send the UI for rendering
// 	return s
// }
