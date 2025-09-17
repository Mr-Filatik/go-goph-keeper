// Package view содержит логику для работы с пользовательским интерфейсом.
package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// IScreen определяет основные функции для всех экранов приложения.
type IScreen interface {
	// fmt.Stringer выводит экран в виде строки для отрисовки.
	fmt.Stringer

	// GetHints возвращает список подсказок для данного экрана.
	GetHints() []Hint

	// Update описывает все действия для данного окна.
	Update(msg tea.Msg) (*MainModel, tea.Cmd)

	// ValidateScreenData проверяет и корректирует данные для текущего экрана.
	ValidateScreenData()
}

// Hint описывает подсказки для пользователя.
type Hint struct {
	actionName string   // Название действия.
	buttons    []string // Кнопки.
}

type user struct {
	Login string
}
