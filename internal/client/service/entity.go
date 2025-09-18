// Package service содержит общие данные для всех сервисов с логикой приложения.
package service

import "time"

// ItemType описывает тип хранимой информации.
type ItemType string

const (
	// ItemLogin - пара логин/пароль.
	PasswordTypeLogin ItemType = "login"

	// ItemText - произвольный текст.
	PasswordTypeText ItemType = "text"

	// ItemBinary - произвольный бинарь.
	PasswordTypeBinary ItemType = "binary"

	// ItemCard - данные банковских карт.
	PasswordTypeCard ItemType = "card"
)

// Password описывает сущность пароля.
type Password struct {
	ID          string            `json:"id"`
	Type        ItemType          `json:"type"`
	Title       string            `json:"title"`
	Description string            `json:"desc"`
	Meta        map[string]string `json:"meta"` // произвольная метаинфа
	Login       string            `json:"username"`
	Password    string            `json:"data"` // шифротекст (opaque)
	Version     int64             `json:"version"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}
