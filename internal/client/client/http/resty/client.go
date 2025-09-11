// Package resty предоставляет функционал для работы с клиентом на основе github.com/go-resty/resty/v2.
package resty

import (
	"context"
	"net/http"

	restylib "github.com/go-resty/resty/v2"
	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
)

// Client - клиент для отправки запросов к серверу.
type Client struct {
	restyClient   *restylib.Client
	log           logger.Logger
	serverAddress string
}

// ClientConfig - структура, содержащая основные параметры для Client.
type ClientConfig struct {
	ServerAddress string
}

// NewClient создаёт новый экземпляр *Client.
func NewClient(config *ClientConfig, l logger.Logger) *Client {
	client := &Client{
		serverAddress: config.ServerAddress,
		log:           l,
		restyClient:   nil,
	}

	return client
}

// Start запускает экземпляр Client.
func (c *Client) Start(_ context.Context) error {
	c.log.Info(
		"Start Client...",
		"address", c.serverAddress,
	)

	c.restyClient = restylib.New()

	c.log.Info("Start Client is successful")

	return nil
}

// Shutdown мягко завершает работу Client.
func (c *Client) Shutdown(_ context.Context) error {
	c.log.Info("Client shutdown starting...")

	transport, isTransport := c.restyClient.GetClient().Transport.(*http.Transport)
	if isTransport {
		transport.CloseIdleConnections()
	}

	c.log.Info("Client shutdown is successful")

	return nil
}

// Close завершает работу Client.
func (c *Client) Close() error {
	c.log.Info("Client close not implemented")

	return nil
}
