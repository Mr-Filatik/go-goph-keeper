package file_test

import (
	"io/fs"
	"testing"
	"testing/fstest"
	"time"

	"github.com/mr-filatik/go-goph-keeper/internal/server/file"
	"github.com/stretchr/testify/assert"
)

/*
	===== NewClientFileStorage =====
*/

func TestNewClientFileStorage(t *testing.T) {
	t.Parallel()

	testFS := fstest.MapFS{}

	strorage := file.NewClientFileStorage(testFS)

	assert.NotEmpty(t, strorage)
}

/*
	===== GetFileInfo =====
*/

type argGetFileInfo struct {
	os string
}

type wantGetFileInfo struct {
	path string
	name string
	err  error
}

type testGetFileInfo struct {
	name string
	args argGetFileInfo
	want wantGetFileInfo
}

func createTestsForGetFileInfo() []testGetFileInfo {
	tests := []testGetFileInfo{
		{
			name: "linux client",
			args: argGetFileInfo{
				os: "linux",
			},
			want: wantGetFileInfo{
				path: "static/",
				name: "client-linux.exe",
				err:  nil,
			},
		},
		{
			name: "macos client",
			args: argGetFileInfo{
				os: "macos",
			},
			want: wantGetFileInfo{
				path: "static/",
				name: "client-macos.exe",
				err:  nil,
			},
		},
		{
			name: "windows client",
			args: argGetFileInfo{
				os: "windows",
			},
			want: wantGetFileInfo{
				path: "static/",
				name: "client-windows.exe",
				err:  nil,
			},
		},
		{
			name: "unknown client",
			args: argGetFileInfo{
				os: "unknown",
			},
			want: wantGetFileInfo{
				path: "",
				name: "",
				err:  file.ErrUncorrectClientOS,
			},
		},
	}

	return tests
}

func TestGetFileInfo(t *testing.T) {
	t.Parallel()

	tests := createTestsForGetFileInfo()

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

	storage := file.NewClientFileStorage(testFS)

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			path, name, err := storage.GetFileInfo(internalTest.args.os)

			assert.Equal(t, internalTest.want.path, path)
			assert.Equal(t, internalTest.want.name, name)

			if internalTest.want.err != nil {
				assert.ErrorIs(t, err, internalTest.want.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

/*
	===== GetFileData =====
*/

type argGetFileData struct {
	path string
	name string
}

type wantGetFileData struct {
	data []byte
	err  error
}

type testGetFileData struct {
	name string
	args argGetFileData
	want wantGetFileData
}

func createTestsForGetFileData() []testGetFileData {
	tests := []testGetFileData{
		{
			name: "linux client",
			args: argGetFileData{
				path: "static/",
				name: "client-linux.exe",
			},
			want: wantGetFileData{
				data: []byte("test-linux-binary"),
				err:  nil,
			},
		},
		{
			name: "macos client",
			args: argGetFileData{
				path: "static/",
				name: "client-macos.exe",
			},
			want: wantGetFileData{
				data: []byte("test-macos-binary"),
				err:  nil,
			},
		},
		{
			name: "windows client",
			args: argGetFileData{
				path: "static/",
				name: "client-windows.exe",
			},
			want: wantGetFileData{
				data: []byte("test-windows-binary"),
				err:  nil,
			},
		},
		{
			name: "unknown client",
			args: argGetFileData{
				path: "static/",
				name: "client-unknown.exe",
			},
			want: wantGetFileData{
				data: nil,
				err:  fs.ErrNotExist,
			},
		},
	}

	return tests
}

func TestGetFileData(t *testing.T) {
	t.Parallel()

	tests := createTestsForGetFileData()

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

	storage := file.NewClientFileStorage(testFS)

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			data, err := storage.GetFileData(internalTest.args.path, internalTest.args.name)

			assert.Equal(t, internalTest.want.data, data)

			if internalTest.want.err != nil {
				assert.ErrorIs(t, err, internalTest.want.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
