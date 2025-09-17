// Package client предоставляет функционал для обработчиков запросов для скачивания клиента.
package client

// InfoResp описывает ответ для запроса информации о доступных клиентах.
type InfoResp struct {
	Path    string   `json:"path"`
	Example string   `json:"example"`
	OS      []string `json:"os"`
}
