package moduleconfig

import (
	"path"
)

type FileSystem interface {
	FileExists(path string) (bool, error)
}

type ModuleConfigService struct {
	fileSystem FileSystem
}

func NewModuleConfigService(fileSystemUtil FileSystem) *ModuleConfigService {
	return &ModuleConfigService{
		fileSystem: fileSystemUtil,
	}
}

func (s *ModuleConfigService) PreventOverwrite(directory, fileName string, overwrite bool) error {
	exists, err := s.fileSystem.FileExists(path.Join(directory, fileName))
	if err != nil {
		return err
	}

	if exists && !overwrite {
		return ErrFileExists
	}

	return nil
}
