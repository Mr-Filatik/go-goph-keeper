// Package main предоставляет точку входа для приложения сервера.
package main

import (
	"github.com/mr-filatik/go-goph-keeper/internal/server"
)

// main является точкой входа для запуска приложения сервера.
func main() {
	server.Run()
}
