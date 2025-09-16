// Package service содержит общие данные для всех сервисов с логикой приложения.
package service

// Password описывает сущность пароля.
type Password struct {
	ID          string
	Title       string
	Description string
	Login       string
	Password    string
	Notes       string
}
