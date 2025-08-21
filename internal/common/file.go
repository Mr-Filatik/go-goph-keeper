// Package common предоставляет общий функционал для приложений.
package common

import (
	"errors"
	"fmt"
	"io/fs"
)

// ExistsFS - проверяет существует ли указанный файл.
func ExistsFS(fsys fs.FS, path string) (bool, error) {
	_, err := fs.Stat(fsys, path)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, fmt.Errorf("unexpected error: %w", err)
}
