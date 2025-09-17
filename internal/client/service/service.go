// Package service содержит общие данные для всех сервисов с логикой приложения.
package service

import "context"

// IService - интерфейс для основной логики приложения.
type IService interface {
	Login(ctx context.Context, login, password string) error
	Register(ctx context.Context, login, password string) error
	Logout(ctx context.Context) error

	GetPasswords(ctx context.Context) ([]Password, error)
	GetPassword(ctx context.Context, passID string) (string, error)
	AddPassword(ctx context.Context, pass Password) (string, error)
	ChangePassword(ctx context.Context, pass Password) error
	RemovePassword(ctx context.Context, passID string) error
}
