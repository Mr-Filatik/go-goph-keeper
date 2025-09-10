// Package file предоставляет функционал для работы с файловой системой.
package file

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"

	root "github.com/mr-filatik/go-goph-keeper"
	"github.com/mr-filatik/go-goph-keeper/internal/common"
)

var (
	// ErrFileNotFound указывает на то, что файл не найден.
	ErrFileNotFound = errors.New("client file not found")

	// ErrUncorrectClientOS указывает на то, что указанная OC неправильная.
	ErrUncorrectClientOS = errors.New("uncorrect client os")
)

// IFileStorage общий интерфейс для всех файловых хранилищ.
type IFileStorage interface {
	GetFileInfo(id string) (path string, name string, err error)
	GetFileData(path, name string) ([]byte, error)
}

// ClientFileStorage структура для доступа к бинарникам клиентов.
type ClientFileStorage struct {
	stor fs.FS
}

// NewClientFileStorage создаёт новый экземпляр ClientFileStorage.
func NewClientFileStorage(stor fs.FS) *ClientFileStorage {
	storage := &ClientFileStorage{
		stor: stor,
	}

	return storage
}

const (
	// ClientOSLinux описывает операционную систему linux.
	ClientOSLinux = "linux"

	// ClientOSMacOS описывает операционную систему macos.
	ClientOSMacOS = "macos"

	// ClientOSWindows описывает операционную систему windows.
	ClientOSWindows = "windows"
)

// GetFileInfo выводит информацию о файле по его уникальному идентификатору.
// Возвращает путь до файла, название файла с расширением и ошибку.
//
// Параметры:
//   - id: уникальный идентификатор файла.
func (s *ClientFileStorage) GetFileInfo(id string) (string, string, error) {
	var name string

	switch id {
	case ClientOSWindows:
		name = "client-windows.exe"
	case ClientOSMacOS:
		name = "client-macos.exe"
	case ClientOSLinux:
		name = "client-linux.exe"
	default:
		return "", "", ErrUncorrectClientOS
	}

	path := strings.Join([]string{root.DirStatic, "/"}, "")

	fullpath := strings.Join([]string{path, name}, "")

	ok, existErr := common.ExistsFS(root.EmbedStatic, fullpath)
	if !ok {
		if existErr != nil {
			return "", "", ErrFileNotFound
		}

		return "", "", ErrFileNotFound
	}

	return path, name, nil
}

// GetFileData выводит содержимое файла.
//
// Параметры:
//   - path: путь до файла;
//   - name: имя файла с расширением.
func (s *ClientFileStorage) GetFileData(path, name string) ([]byte, error) {
	fullpath := strings.Join([]string{path, name}, "")

	file, openErr := s.stor.Open(fullpath)

	if openErr != nil {
		return nil, fmt.Errorf("open file: %w", openErr)
	}

	closeErr := file.Close()
	if closeErr != nil {
		return nil, fmt.Errorf("close file: %w", closeErr)
	}

	data, readErr := io.ReadAll(file)
	if readErr != nil {
		return nil, fmt.Errorf("read bytes from file: %w", readErr)
	}

	return data, nil
}
