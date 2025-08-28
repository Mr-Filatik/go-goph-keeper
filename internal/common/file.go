// Package common предоставляет общий функционал для приложений.
package common

import (
	"errors"
	"fmt"
	"io/fs"
)

var errArgumentIsEmpty = errors.New("argument is empty")

// GetErrArgumentIsEmpty - возвращает ошибку 'argument is empty'.
func GetErrArgumentIsEmpty() error {
	return errArgumentIsEmpty
}

// ExistsFS - проверяет существует ли указанный файл.
func ExistsFS(fsys fs.FS, path string) (bool, error) {
	if path == "" {
		return false, fmt.Errorf("path: %w", errArgumentIsEmpty)
	}

	_, err := fs.Stat(fsys, path)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, fmt.Errorf("unexpected error: %w", err)
}
