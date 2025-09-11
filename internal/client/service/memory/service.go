package memory

import (
	"context"

	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
)

// Service - клиент для отправки запросов к серверу.
type Service struct {
	log logger.Logger
}

// NewService создаёт новый экземпляр *Service.
func NewService(l logger.Logger) *Service {
	client := &Service{
		log: l,
	}

	return client
}

func (s *Service) Ping(_ context.Context) error {
	return nil
}

func (s *Service) Login(_ context.Context, login, password string) error {
	return nil
}

func (s *Service) Register(_ context.Context, login, password string) error {
	return nil
}

func (s *Service) Logout(_ context.Context) error {
	return nil
}

func (s *Service) GetPasswords(_ context.Context) ([]string, error) {
	return []string{}, nil
}

func (s *Service) GetPassword(_ context.Context, passID string) (string, error) {
	return "", nil
}

// add password
// change password
// remove password
