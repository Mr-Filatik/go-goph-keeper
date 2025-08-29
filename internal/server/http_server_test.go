// Package server_test предоставляет функционал для тестирования серверов.
package server_test

import (
	"context"
	"testing"

	"github.com/mr-filatik/go-goph-keeper/internal/server"
	"github.com/mr-filatik/go-goph-keeper/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPServer(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()

	conf := &server.HTTPServerConfig{
		Address: "127.0.0.1:0",
	}
	serv := server.NewHTTPServer(conf, nil, mockLogger)

	assert.NotEmpty(t, serv)
}

func TestHTTPServer_Start(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()

	conf := &server.HTTPServerConfig{
		Address: "127.0.0.1:0",
	}
	serv := server.NewHTTPServer(conf, nil, mockLogger)

	ctx := context.Background()
	err := serv.Start(ctx)

	assert.NoError(t, err)
}

func TestHTTPServer_Shutdown(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()

	conf := &server.HTTPServerConfig{
		Address: "127.0.0.1:0",
	}
	serv := server.NewHTTPServer(conf, nil, mockLogger)

	ctx := context.Background()
	err := serv.Start(ctx)

	require.NoError(t, err)

	err = serv.Shutdown(ctx)

	require.NoError(t, err)
}

func TestHTTPServer_Close(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()

	conf := &server.HTTPServerConfig{
		Address: "127.0.0.1:0",
	}
	serv := server.NewHTTPServer(conf, nil, mockLogger)

	ctx := context.Background()
	err := serv.Start(ctx)

	require.NoError(t, err)

	err = serv.Shutdown(ctx)

	require.NoError(t, err)

	err = serv.Close()

	require.NoError(t, err)
}
