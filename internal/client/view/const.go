// Package view содержит логику для работы с пользовательским интерфейсом.
package view

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
)

// зарефакторить каким-то образом работу с шагами алгоритмов.
const (
	stepInit = iota
	stepOne
	stepTwo
)
