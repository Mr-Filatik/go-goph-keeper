// Package server предоставляет функционал для запуска приложения сервера.
package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mr-filatik/go-goph-keeper/internal/common/logger"
)

const (
	timeoutIdle       = 5 * time.Second
	timeoutRead       = 5 * time.Second
	timeoutReadHeader = 5 * time.Second
	timeoutWrite      = 10 * time.Second
)

// HTTPServer представляет HTTP-сервер приложения.
type HTTPServer struct {
	server  *http.Server  // сервер
	router  *chi.Mux      // роутер
	log     logger.Logger // логгер
	address string        // адрес сервера
}

// HTTPServerConfig - конфиг для создания HTTPServer.
type HTTPServerConfig struct {
	Address string
}

// NewHTTPServer создаёт и инициализирует новый экзепляр *HTTPServer.
//
// Параметры:
//   - conf: конфиг сервера;
//   - log: логгер.
func NewHTTPServer(conf *HTTPServerConfig, log logger.Logger) *HTTPServer {
	log.Info("HTTPServer creating...")

	srv := &HTTPServer{
		server:  nil,
		address: conf.Address,
		router:  chi.NewRouter(),
		log:     log,
	}

	log.Info("HTTPServer create is successful")

	return srv
}

// Start запускает экземпляр HTTPServer.
func (s *HTTPServer) Start(ctx context.Context) error {
	s.log.Info(
		"HTTPServer starting...",
		"address", s.address,
	)

	tslNextProto := make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)

	s.server = &http.Server{
		Addr: s.address,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		ConnContext:                  nil,
		ConnState:                    nil,
		DisableGeneralOptionsHandler: false,
		ErrorLog:                     nil,
		Handler:                      nil,
		IdleTimeout:                  timeoutIdle,
		MaxHeaderBytes:               http.DefaultMaxHeaderBytes,
		ReadHeaderTimeout:            timeoutReadHeader,
		ReadTimeout:                  timeoutRead,
		TLSConfig:                    nil,
		TLSNextProto:                 tslNextProto,
		WriteTimeout:                 timeoutWrite,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("Error in HTTPServer", err)
		}
	}()

	s.log.Info("HTTPServer start is successful")

	return nil
}

// Shutdown мягко завершает работу HTTPServer.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	s.log.Info("HTTPServer shutdown starting...")

	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("HTTPServer shutdown error: %w", err)
	}

	s.log.Info("HTTPServer shutdown is successful")

	return nil
}

// Close завершает работу HTTPServer.
func (s *HTTPServer) Close() error {
	s.log.Info("HTTPServer close starting...")

	err := s.server.Close()
	if err != nil {
		return fmt.Errorf("HTTPServer close error: %w", err)
	}

	s.log.Info("HTTPServer close is successful")

	return nil
}
