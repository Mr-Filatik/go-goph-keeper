package view

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginScreen struct {
	mainModel     *teaModel
	LoginInput    textinput.Model
	PasswordInput textinput.Model
	ErrMessage    string
}

func NewLoginScreen(mod *teaModel) *LoginScreen {
	// login input
	loginInput := textinput.New()
	loginInput.Placeholder = "your email"
	loginInput.CharLimit = 64
	loginInput.Focus()

	// password inputs
	passInput := textinput.New()
	passInput.Placeholder = "your password"
	passInput.CharLimit = 64
	passInput.EchoMode = textinput.EchoPassword
	passInput.EchoCharacter = '•'

	return &LoginScreen{
		mainModel:     mod,
		LoginInput:    loginInput,
		PasswordInput: passInput,
		ErrMessage:    "",
	}
}

func (s *LoginScreen) LoadScreen(fnc func()) IScreen {
	// login input
	loginInput := textinput.New()
	loginInput.Placeholder = "your email"
	loginInput.CharLimit = 64
	loginInput.Focus()

	// password inputs
	passInput := textinput.New()
	passInput.Placeholder = "your password"
	passInput.CharLimit = 64
	passInput.EchoMode = textinput.EchoPassword
	passInput.EchoCharacter = '•'

	s.LoginInput = loginInput
	s.PasswordInput = passInput
	s.ErrMessage = ""

	if fnc != nil {
		fnc()
	}

	return s
}

func (s *LoginScreen) String() string {
	view := "\n[Login] Enter email and password:\n"

	view += s.LoginInput.View() + "\n"
	view += s.PasswordInput.View() + "\n\n"

	if s.ErrMessage != "" {
		view += "\n[ERROR]: " + s.ErrMessage + "\n"
	}

	return view
}

func (s *LoginScreen) GetHints() []Hint {
	return []Hint{
		{"Login", []string{KeyEnter}},
		{"Switch", []string{KeyTab, KeyDown, KeyUp}},
		{"Back", []string{KeyEscape}},
		{"Quit", []string{KeyQuit}},
	}
}

func (s *LoginScreen) Action(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch v := msg.(type) {
	// case authProgress:
	//     if ls, ok := m.screenCurrent.(*LoadingScreen); ok {
	//         ls.percent = v.percent
	//         ls.status = v.status
	//     }
	//     return m, nil

	// case authOK:
	//     m.currentUser = &v.u
	//     m.screenCurrent = m.screenPassList.LoadScreen(nil)
	//     return m, nil

	// case authFailed:
	//     // вернуться на логин с ошибкой
	//     if ls, ok := m.screenCurrent.(*LoadingScreen); ok && ls.cancel != nil {
	//         // уже закончилась — cancel не нужен, но пусть будет безопасно
	//         ls.cancel()
	//     }
	//     // переоткроем логин-скрин и передадим сообщение
	//     if lscr, ok := m.screenStart.(*LoginScreen); ok {
	//         lscr.ErrMessage = v.err.Error()
	//         s.mainModel.screenCurrent = lscr.LoadScreen(nil)
	//     }
	//     return s.mainModel, nil

	// case authCanceled:
	//     // вернуться на логин без ошибки (или со своей подписью)
	//     if lscr, ok := s.mainModel.screenStart.(*LoginScreen); ok {
	//         lscr.ErrMessage = "Canceled"
	//         s.mainModel.screenCurrent = lscr.LoadScreen(nil)
	//     }
	//     return s.mainModel, nil
	// }

	key, isKey := msg.(tea.KeyMsg)
	if isKey {
		switch key.String() {
		case KeyQuit:
			return s.mainModel, tea.Quit

		case KeyEscape:
			s.mainModel.screenCurrent = s.mainModel.screenStart.LoadScreen(nil)

			return s.mainModel, nil

		case KeyEnter:
			login := s.LoginInput.Value()
			password := s.PasswordInput.Value()

			if login == "" || password == "" {
				s.ErrMessage = "login and password are required"

				return s.mainModel, nil
			}

			// ctx, _ := context.WithCancel(context.Background())

			// s.mainModel.loadingCmd = authCmd(ctx, login, password)
			s.mainModel.screenCurrent = s.mainModel.screenLoading.LoadScreen(func() {
				s.mainModel.screenLoading.title = "Test"
				s.mainModel.screenLoading.desc = "test test test"

				s.mainModel.screenLoading.login = login
				s.mainModel.screenLoading.pass = password
			})

			// s.mainModel.screenCurrent = s.mainModel.screenPassList.LoadScreen(func() {
			// 	s.mainModel.currentUser = &user{
			// 		Login: s.LoginInput.Value(),
			// 	}
			// })

			// s.ErrMessage = "happy (" + s.LoginInput.Value() + ") (" + s.PasswordInput.Value() + ")"

			return s.mainModel, tea.Tick(150*time.Millisecond, func(time.Time) tea.Msg { return LoadingProgressMsg{} })

		case KeyTab, KeyDown, KeyUp:
			if s.LoginInput.Focused() {
				s.LoginInput.Blur()
				s.PasswordInput.Focus()
			} else {
				s.PasswordInput.Blur()
				s.LoginInput.Focus()
			}

			return s.mainModel, nil
		}
	}

	var cmd tea.Cmd
	if s.LoginInput.Focused() {
		s.LoginInput, cmd = s.LoginInput.Update(key)
	} else {
		s.PasswordInput, cmd = s.PasswordInput.Update(key)
	}

	return s.mainModel, cmd
}

// type authOK struct{ u user }
// type authFailed struct{ err error }
// type authCanceled struct{}
// type authProgress struct {
// 	percent float64
// 	status  string
// }

// func (s *LoginScreen) authAction(msg tea.Msg) (tea.Model, tea.Cmd) {

// }

// func authCmd(ctx context.Context, login, pass string) tea.Cmd {
// 	return func() tea.Msg {
// 		// пример «прогресса» (замените на реальную работу)
// 		ticker := time.NewTicker(500 * time.Millisecond)
// 		defer ticker.Stop()
// 		p := 0.0
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return authCanceled{}
// 			case <-ticker.C:
// 				p += 5
// 				if p >= 100 {
// 					// тут сделайте реальную проверку логина/пароля
// 					if login == "demo" && pass == "demo" {
// 						return authOK{u: user{Login: login}}
// 					}
// 					return authFailed{err: fmt.Errorf("invalid credentials")}
// 				}
// 				// можно периодически возвращать прогресс
// 				return tea.Batch(
// 					func() tea.Msg { return authProgress{percent: p, status: "Authorizing…"} },
// 					authCmd(ctx, login, pass), // продолжить
// 				)
// 			}
// 		}
// 	}
// }
