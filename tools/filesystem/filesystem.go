package filesystem

import (
	"errors"
	"os"
)

type FileSystemUtil struct{}

func (u *FileSystemUtil) FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil

	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil

	} else {
		return false, err
	}
}
