package client_test

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"time"

	"github.com/mr-filatik/go-goph-keeper/internal/server/file"
	"github.com/mr-filatik/go-goph-keeper/internal/server/handler"
	"github.com/mr-filatik/go-goph-keeper/internal/server/handler/client"
	"github.com/mr-filatik/go-goph-keeper/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
	===== NewHandler =====
*/

func TestNewHandler(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()

	mainHandler := handler.NewHandler(nil, mockLogger)
	clinetHandler := client.NewHandler(*mainHandler)

	assert.NotEmpty(t, clinetHandler)
}

/*
	===== ClientInfo =====
*/

func TestClientInfo(t *testing.T) {
	t.Parallel()

	mockLogger := testutil.NewMockLogger()

	mainHandler := handler.NewHandler(nil, mockLogger)
	clientHandler := client.NewHandler(*mainHandler)

	req := httptest.NewRequest(http.MethodGet, "/client", http.NoBody)
	recorder := httptest.NewRecorder()

	clientHandler.ClientInfo(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	expected := client.InfoResp{
		Path:    "/client/{os}",
		Example: "http://localhost:8080/client/linux",
		OS: []string{
			"linux",
			"macos",
			"windows",
		},
	}

	var resp client.InfoResp

	decoder := json.NewDecoder(recorder.Body)
	require.NotNil(t, decoder)

	err := decoder.Decode(&resp)
	require.NoError(t, err)

	assert.Equal(t, expected.Path, resp.Path)
	assert.Equal(t, expected.Example, resp.Example)
	assert.ElementsMatch(t, expected.OS, resp.OS)
}

/*
	===== ClientDownload =====
*/

type argClientDownload struct {
	os   string
	path string
}

type wantClientDownload struct {
	contentType string
	statusCode  int
}

type testClientDownload struct {
	name string
	args argClientDownload
	want wantClientDownload
}

func createTestsForClientDownload() []testClientDownload {
	tests := []testClientDownload{
		{
			name: "linux client",
			args: argClientDownload{
				os:   "linux",
				path: "/client/linux",
			},
			want: wantClientDownload{
				statusCode:  200,
				contentType: "application/octet-stream",
			},
		},
		{
			name: "macos client",
			args: argClientDownload{
				os:   "macos",
				path: "/client/macos",
			},
			want: wantClientDownload{
				statusCode:  200,
				contentType: "application/x-mach-binary",
			},
		},
		{
			name: "windows client",
			args: argClientDownload{
				os:   "windows",
				path: "/client/windows",
			},
			want: wantClientDownload{
				statusCode:  200,
				contentType: "application/x-msdownload",
			},
		},
		{
			name: "unknown client",
			args: argClientDownload{
				os:   "unknown",
				path: "/client/unknown",
			},
			want: wantClientDownload{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	return tests
}

func TestClientDownload(t *testing.T) {
	t.Parallel()

	tests := createTestsForClientDownload()

	mockLogger := testutil.NewMockLogger()

	testFS := fstest.MapFS{
		"static/client-linux.exe": &fstest.MapFile{
			Data:    []byte("test-linux-binary"),
			Mode:    fs.ModePerm,
			ModTime: time.Now(),
			Sys:     nil,
		},
		"static/client-macos.exe": &fstest.MapFile{
			Data:    []byte("test-macos-binary"),
			Mode:    fs.ModePerm,
			ModTime: time.Now(),
			Sys:     nil,
		},
		"static/client-windows.exe": &fstest.MapFile{
			Data:    []byte("test-windows-binary"),
			Mode:    fs.ModePerm,
			ModTime: time.Now(),
			Sys:     nil,
		},
	}
	mockFileStorage := file.NewClientFileStorage(testFS) // Переписать под реальный мок

	mainHandler := handler.NewHandler(nil, mockLogger)
	clientHandlerOptions := client.WithCustomFileStorage(mockFileStorage)
	clientHandler := client.NewHandler(*mainHandler, clientHandlerOptions)

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, internalTest.args.path, http.NoBody)
			req.SetPathValue("os", internalTest.args.os)

			recorder := httptest.NewRecorder()

			clientHandler.ClientDownload(recorder, req)

			assert.Equal(t, internalTest.want.statusCode, recorder.Code)
			assert.Equal(t, internalTest.want.contentType, recorder.Header().Get("Content-Type"))
		})
	}
}
