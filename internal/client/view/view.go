package view

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenStart screen = iota
	screenLogin
	screenRegister
	screenLoading
	screenPasswordList
	screenPasswordDetails
)

type item struct {
	Title       string
	Description string
	Username    string
	Password    string // never display; only copy on request
	Notes       string
}

type availableAction struct {
	name    string
	buttons []string
}

type model struct {
	currentScreen screen
	currentUser   string // current user email

	// global/app info always shown in header
	appName      string
	buildVersion string
	buildDate    string
	buildCommit  string

	// login inputs
	loginInput    textinput.Model
	passwordInput textinput.Model
	loginErr      string

	availableActions []availableAction
	indexInput       int // индекс для выборов в меню (один на разные окна/обнуляется при переходе)

	// loading screen data
	loadingInfo    string
	loadingPercent float64

	// data
	items     []item
	listIndex int
	selected  *item
	statusMsg string // ephemeral footer message

	screenCurrent  IScreen
	screenStart    *StartScreen
	screenLogin    *LoginScreen
	screenRegister *RegisterScreen
}

func initialModel() model {
	// login input
	loginInput := textinput.New()
	loginInput.Placeholder = "email or login"
	loginInput.CharLimit = 64
	loginInput.Focus()

	// password inputs
	passInput := textinput.New()
	passInput.Placeholder = "password"
	passInput.CharLimit = 64
	passInput.EchoMode = textinput.EchoPassword
	passInput.EchoCharacter = '•'

	// sample data (replace with your backend later)
	data := []item{
		{"GitHub", "Personal account", "vlad", "ghp_example_password", "2FA: TOTP in Authy"},
		{"GMail", "Work", "vladislav", "gmail_app_password", "App password only"},
		{"AWS", "Prod account", "admin", "supersecretkey", "Use IAM roles"},
	}

	mod := &model{
		currentScreen: screenStart,
		appName:       "AppName",
		buildVersion:  "Version",
		buildDate:     "Build",
		buildCommit:   "Commit",
		currentUser:   "user",
		loginInput:    loginInput,
		passwordInput: passInput,
		items:         data,
		listIndex:     0,
	}

	mod.screenStart = NewStartScreen(mod)
	mod.screenLogin = NewLoginScreen(mod)
	mod.screenRegister = NewRegisterScreen(mod)

	mod.screenCurrent = mod.screenStart

	return *mod
}

func (m model) Init() tea.Cmd { return nil }

const (
	KeyUp     = "up"
	KeyUpWASD = "w"

	KeyDown     = "down"
	KeyDownWASD = "s"

	KeyNext = "n"
	KeyPrev = "p"

	KeyTab = "tab"

	KeyEnter = "enter"

	KeyEscape      = "esc"
	KeyEscapeShort = "b"

	KeyQuit      = "ctrl+c"
	KeyQuitShort = "q"

	KeyCopy = "c"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.screenCurrent != nil {
			return m.screenCurrent.Action(msg)
		}

		// switch m.currentScreen {
		// case screenLogin:
		// 	switch msg.String() {
		// 	case KeyQuit, KeyQuitShort, KeyEscape:
		// 		return m, tea.Quit
		// 	case KeyTab:
		// 		if m.loginInput.Focused() {
		// 			m.loginInput.Blur()
		// 			m.passwordInput.Focus()
		// 		} else {
		// 			m.passwordInput.Blur()
		// 			m.loginInput.Focus()
		// 		}
		// 		return m, nil
		// 	case KeyEnter:
		// 		// TODO: replace with real auth
		// 		if m.loginInput.Value() == "" || m.passwordInput.Value() == "" {
		// 			m.loginErr = "login and password are required"
		// 			return m, nil
		// 		}
		// 		m.currentScreen = screenPasswordList
		// 		m.statusMsg = "Logged in"
		// 		m.loginErr = ""
		// 		return m, nil
		// 	}
		// 	// delegate typing to inputs
		// 	var cmd tea.Cmd
		// 	if m.loginInput.Focused() {
		// 		m.loginInput, cmd = m.loginInput.Update(msg)
		// 	} else {
		// 		m.passwordInput, cmd = m.passwordInput.Update(msg)
		// 	}
		// 	return m, cmd

		// case screenPasswordList:
		// 	switch msg.String() {
		// 	case KeyQuit, KeyQuitShort:
		// 		return m, tea.Quit
		// 	case KeyUp, KeyUpWASD:
		// 		if m.listIndex > 0 {
		// 			m.listIndex--
		// 		}
		// 		return m, nil
		// 	case KeyDown, KeyDownWASD:
		// 		if m.listIndex < len(m.items)-1 {
		// 			m.listIndex++
		// 		}
		// 		return m, nil
		// 	case KeyEnter:
		// 		m.selected = &m.items[m.listIndex]
		// 		m.currentScreen = screenPasswordDetails
		// 		m.statusMsg = "Opened details"
		// 		return m, nil
		// 	}

		// case screenPasswordDetails:
		// 	switch msg.String() {
		// 	case KeyQuit, KeyQuitShort:
		// 		return m, tea.Quit
		// 	case KeyEscape, KeyEscapeShort:
		// 		m.currentScreen = screenPasswordList
		// 		m.selected = nil
		// 		m.statusMsg = "Back to list"
		// 		return m, nil
		// 	case KeyCopy:
		// 		if m.selected != nil {
		// 			_ = clipboard.WriteAll(m.selected.Password)
		// 			m.statusMsg = "Password copied to clipboard"
		// 		}
		// 		return m, nil
		// 	}
		// }
	}

	return m, nil
}

func (m model) View() string {
	view := m.header() // вынести подсказки как header [Enter] Select [Tab/Up/Down/N/P/W/S] Switch [Esc/Q/B/Ctrl+C] Quit\n\n
	// записывать их в отдельную переменную (действие и []кнопок)

	if m.screenCurrent != nil {
		view += m.screenCurrent.String()
	}

	// switch m.currentScreen {
	// case screenStart:
	// 	view += m.screenStart.String() // переделать на текущее окно
	// case screenLogin:
	// 	view += m.viewLoginScreen()
	// case screenPasswordList:
	// 	view += m.viewPasswordListScreen()
	// case screenPasswordDetails:
	// 	view += m.viewPasswordDetailsScreen()
	// }

	view += "\n" + m.footer() + "\n"

	return view
}

func (m model) viewLoginScreen() string {
	view := "Login\n"
	view += m.loginInput.View() + "\n"
	view += m.passwordInput.View() + "\n\n"
	view += "[Enter] Sign in [Tab] Switch field [Esc/Ctrl+C] Quit\n"

	if m.loginErr != "" {
		view += "\nERROR: " + m.loginErr + "\n"
	}

	return view
}

func (m model) viewRegisterScreen() string {
	return ""
}

func (m model) viewPasswordListScreen() string {
	view := "Passwords (Up/Down, Enter=Open, q=Quit)\n\n"

	for i, item := range m.items {
		cursor := " "
		if i == m.listIndex {
			cursor = ">"
		}

		view += fmt.Sprintf("%s %s — %s\n", cursor, item.Title, item.Description)
	}

	return view
}

func (m model) viewPasswordDetailsScreen() string {
	var view string

	if m.selected != nil {
		item := *m.selected
		view := fmt.Sprintf("%s\n%s\n\n", item.Title, item.Description)
		view += fmt.Sprintf("Username: %s\n", item.Username)
		view += "Password: [hidden] (press 'c' to copy)\n"

		if item.Notes != "" {
			view += "Notes: " + item.Notes + "\n"
		}

		view += "\n[b] Back [c] Copy password [q] Quit\n"
	}

	return view
}

const (
	lineWidth  = 60
	lineSymbol = "─"
)

func (m model) header() string {
	parts := []string{
		addLine(),
		m.appName,
	}

	if m.currentUser != "" {
		parts = append(parts, "User: "+m.currentUser)
	}

	parts = append(parts, addLine(), "Hints:")

	if m.screenCurrent != nil {
		hints := m.screenCurrent.GetHints()
		if len(hints) != 0 {
			for _, hint := range hints {
				parts = append(parts, hint.actionName+" ["+strings.Join(hint.buttons, "/")+"]")
			}

			parts = append(parts, addLine()+"\n")
		}
	}

	return strings.Join(parts, "\n")
}

func (m model) footer() string {
	return strings.Join([]string{
		addLine(),
		"Build Version: " + m.buildVersion,
		"Build Date: " + m.buildDate,
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
