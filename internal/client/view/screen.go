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
	Action(key tea.KeyMsg) (tea.Model, tea.Cmd)
}

// Hint описывает подсказки для пользователя.
type Hint struct {
	actionName string   // Название действия.
	buttons    []string // Кнопки.
}
