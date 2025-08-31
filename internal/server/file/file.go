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
	ErrFileNotFound      = errors.New("client file not found")
	ErrUncorrectClientOS = errors.New("uncorrect client os")
)

type IFileStorage interface {
	GetFileInfo(id string) (path string, name string, err error)
	GetFileData(path, name string) ([]byte, error)
}

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
	ClientOSLinux   = "linux"
	ClientOSMacOS   = "macos"
	ClientOSWindows = "windows"
)

func (s *ClientFileStorage) GetFileInfo(id string) (path string, name string, err error) {
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

	path = strings.Join([]string{root.DirStatic, "/"}, "")

	fullpath := strings.Join([]string{path, name}, "")

	ok, existErr := common.ExistsFS(root.EmbedStatic, fullpath)
	if !ok {
		if existErr != nil {
			return "", "", ErrFileNotFound
		}

		return "", "", ErrFileNotFound
	}

	return
}

func (s *ClientFileStorage) GetFileData(path, name string) ([]byte, error) {
	fullpath := strings.Join([]string{path, name}, "")

	file, openError := s.stor.Open(fullpath)
	if openError != nil {
		return nil, fmt.Errorf("open file: %w", openError)
	}

	var buff []byte
	_, readErr := io.ReadFull(file, buff)
	if readErr != nil {
		return nil, fmt.Errorf("read bytes from file: %w", readErr)
	}

	return buff, nil
}
