package memory

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mr-filatik/go-goph-keeper/internal/client/service"
	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
)

// Service - клиент для отправки запросов к серверу.
type Service struct {
	log   logger.Logger
	token string
	pass  []service.Password
}

// NewService создаёт новый экземпляр *Service.
func NewService(l logger.Logger) *Service {
	client := &Service{
		log:   l,
		token: "",
		pass: []service.Password{
			{
				ID:          "1",
				Title:       "Email",
				Description: "Личная почта Gmail",
				Login:       "demo@gmail.com",
				Password:    "p@ssw0rd123",
				Type:        service.PasswordTypeLogin,
			},
			{
				ID:          "2",
				Title:       "GitHub",
				Description: "Аккаунт разработчика",
				Login:       "devuser",
				Password:    "gh_token_abc123",
				Type:        service.PasswordTypeLogin,
			},
			{
				ID:          "3",
				Title:       "Bank",
				Description: "Интернет-банк",
				Login:       "client1234",
				Password:    "B@nkSecure!",
				Type:        service.PasswordTypeLogin,
			},
			{
				ID:          "4",
				Title:       "AWS",
				Description: "Amazon Web Services root",
				Login:       "root@company.com",
				Password:    "aws-secret-key",
				Type:        service.PasswordTypeLogin,
			},
			{
				ID:          "5",
				Title:       "Spotify",
				Description: "Музыка",
				Login:       "musiclover",
				Password:    "spotify!234",
				Type:        service.PasswordTypeLogin,
			},
			{
				ID:          "6",
				Title:       "Telegram",
				Description: "Мессенджер",
				Login:       "+1234567890",
				Password:    "tg_pass_789",
				Type:        service.PasswordTypeLogin,
			},
			{
				ID:          "7",
				Title:       "WorkMail",
				Description: "Корпоративная почта",
				Login:       "user@company.com",
				Password:    "C0rpPass!",
				Type:        service.PasswordTypeLogin,
			},
			{
				ID:          "8",
				Title:       "VPN",
				Description: "Доступ в корпоративную сеть",
				Login:       "vpnuser",
				Password:    "vpn-strong-key",
				Type:        service.PasswordTypeLogin,
			},
			{
				ID:          "9",
				Title:       "Facebook",
				Description: "Личный аккаунт",
				Login:       "fb.demo",
				Password:    "fb!secure456",
				Type:        service.PasswordTypeLogin,
			},
			{
				ID:          "10",
				Title:       "DockerHub",
				Description: "Образы контейнеров",
				Login:       "dockuser",
				Password:    "d0ckerHUB!",
				Type:        service.PasswordTypeLogin,
			},
		},
	}

	return client
}

func (s *Service) Login(ctx context.Context, login, password string) error {
	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return context.Canceled

	case <-timer.C:
		if login != "demo" || password != "demo" {
			return fmt.Errorf("invalid credentials: %w", errors.New("login or password"))
		}

		s.token = "1"

		return nil
	}
}

func (s *Service) Register(ctx context.Context, login, password string) error {
	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return context.Canceled

	case <-timer.C:
		if login != "demo" || password != "demo" {
			return fmt.Errorf("user already register: %w", errors.New("login  found"))
		}

		s.token = "2"

		return nil
	}
}

func (s *Service) Logout(ctx context.Context) error {
	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return context.Canceled

	case <-timer.C:
		if s.token == "2" {
			return fmt.Errorf("server not connected: %w", errors.New("ups"))
		}

		s.token = ""

		return nil
	}
}

func (s *Service) GetPasswords(ctx context.Context) ([]service.Password, error) {
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	case <-timer.C:
		// Вернём фейковые данные
		return s.pass, nil
	}
}

func (s *Service) GetPassword(ctx context.Context, passID string) (string, error) {
	timer := time.NewTimer(300 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-timer.C:
		if passID == "" {
			return "", fmt.Errorf("password ID is empty")
		}
		// Заглушка: возвращаем «секрет»
		return "secret-password-for-" + passID, nil
	}
}

func (s *Service) AddPassword(ctx context.Context, pass service.Password) (string, error) {
	timer := time.NewTimer(700 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-timer.C:
		if pass.Title == "" {
			return "", fmt.Errorf("title is required")
		}
		// Фейковый ID
		return "new-id-123", nil
	}
}

func (s *Service) ChangePassword(ctx context.Context, pass service.Password) error {
	timer := time.NewTimer(400 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		if pass.ID == "" {
			return fmt.Errorf("password ID is required")
		}
		// Заглушка: просто успех
		return nil
	}
}

func (s *Service) RemovePassword(ctx context.Context, passID string) error {
	timer := time.NewTimer(300 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-timer.C:
		if passID == "2" {
			return fmt.Errorf("password ID is required")
		}

		return nil
	}
}
