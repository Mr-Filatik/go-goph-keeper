package common_test

import (
	"io/fs"
	"testing"
	"testing/fstest"
	"time"

	"github.com/mr-filatik/go-goph-keeper/internal/common"
	"github.com/stretchr/testify/assert"
)

type argsExistsFS struct {
	path string
}

type wantExistsFS struct {
	err   error
	exist bool
}

type testExistsFS struct {
	name string
	args argsExistsFS
	want wantExistsFS
}

func TestExistsFS(t *testing.T) {
	t.Parallel()

	mockFS := fstest.MapFS{
		"file.txt": &fstest.MapFile{
			Data:    []byte("hello"),
			Mode:    fs.ModeTemporary,
			ModTime: time.Now(),
			Sys:     nil,
		},
		"dir/subfile.txt": &fstest.MapFile{
			Data:    []byte("world"),
			Mode:    fs.ModeTemporary,
			ModTime: time.Now(),
			Sys:     nil,
		},
	}

	tests := createTestsForExistsFS()

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			got, err := common.ExistsFS(mockFS, internalTest.args.path)

			assert.Equal(t, internalTest.want.exist, got)

			if internalTest.want.err != nil {
				assert.ErrorIs(t, err, common.GetErrArgumentIsEmpty())
			}
		})
	}
}

func createTestsForExistsFS() []testExistsFS {
	return []testExistsFS{
		{
			name: "existing file in root directory",
			args: argsExistsFS{
				path: "file.txt",
			},
			want: wantExistsFS{
				exist: true,
				err:   nil,
			},
		},
		{
			name: "existing file in subdirectory",
			args: argsExistsFS{
				path: "dir/subfile.txt",
			},
			want: wantExistsFS{
				exist: true,
				err:   nil,
			},
		},
		{
			name: "non-existent file in root directory",
			args: argsExistsFS{
				path: "missing.txt",
			},
			want: wantExistsFS{
				exist: false,
				err:   nil,
			},
		},
		{
			name: "non-existent file in subdirectory",
			args: argsExistsFS{
				path: "dir/missing.txt",
			},
			want: wantExistsFS{
				exist: false,
				err:   nil,
			},
		},
		{
			name: "empty file path",
			args: argsExistsFS{
				path: "",
			},
			want: wantExistsFS{
				exist: false,
				err:   common.GetErrArgumentIsEmpty(),
			},
		},
	}
}
