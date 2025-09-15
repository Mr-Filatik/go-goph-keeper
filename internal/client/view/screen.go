package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type IScreen interface {
	// fmt.Stringer выводит экран в виде строки для отрисовки.
	fmt.Stringer

	// GetHints возвращает список подсказок для данного экрана.
	GetHints() []Hint

	// Action описывает все действия для данного окна.
	Action(msg tea.Msg) (tea.Model, tea.Cmd)

	// LoadAction описывает действия, необходимые при загрузке окна.
	// Сама функция выполняет действия по обнулению необходимых данных.
	// Через func можно дополнительно установить нужные параметры.
	// Подумать как переделать, чтобы функции были более конкретные.
	// И как сделать функцию необязательной.
	LoadScreen(fnc func())
}

// Hint описывает подсказки для пользователя.
type Hint struct {
	actionName string   // Название действия.
	buttons    []string // Кнопки.
}

type item struct {
	Title       string
	Description string
	Username    string
	Password    string // never display; only copy on request
	Notes       string
}

type user struct {
	Login string
}
