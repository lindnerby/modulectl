package scaffold

import (
	"fmt"
	"path"

	"github.com/kyma-project/modulectl/tools/io"
)

type ModuleConfigService interface {
	PreventOverwrite(directory, moduleConfigFileName string, overwrite bool) error
}

type ManifestService interface {
	GetDefaultManifestContent() string
	WriteManifestFile(content, path string) error
}

type DefaultCRService interface {
	GenerateDefaultCRFile(out io.Out, path string) error
}

type FileSystem interface {
	FileExists(path string) (bool, error)
}

type ScaffoldService struct {
	moduleConfigService ModuleConfigService
	manifestService     ManifestService
	defaultCRService    DefaultCRService
	filesystem          FileSystem
}

func NewScaffoldService(moduleConfigService ModuleConfigService,
	manifestService ManifestService,
	defaultCRService DefaultCRService,
	fileSystem FileSystem) *ScaffoldService {
	return &ScaffoldService{
		moduleConfigService: moduleConfigService,
		manifestService:     manifestService,
		defaultCRService:    defaultCRService,
		filesystem:          fileSystem,
	}
}

func (s *ScaffoldService) CreateScaffold(opts Options) error {
	if err := opts.validate(); err != nil {
		return err
	}

	if err := s.moduleConfigService.PreventOverwrite(opts.Directory, opts.ModuleConfigFileName, opts.ModuleConfigFileOverwrite); err != nil {
		return err
	}

	manifestFilePath := path.Join(opts.Directory, opts.ManifestFileName)
	manifestFileExists, err := s.filesystem.FileExists(manifestFilePath)
	if err != nil {
		return err
	}
	if manifestFileExists {
		opts.Out.Write(fmt.Sprintf("The Manifest file already exists, reusing: %s\n", manifestFilePath))
	} else {
		if err := s.manifestService.WriteManifestFile(s.manifestService.GetDefaultManifestContent(), manifestFilePath); err != nil {
			return err
		}
		opts.Out.Write(fmt.Sprintf("Generated a blank Manifest file: %s\n", manifestFilePath))
	}

	defaultCRFilePath := ""
	if opts.defaultCRFileNameConfigured() {
		defaultCRFilePath = path.Join(opts.Directory, opts.DefaultCRFileName)
		if err := s.defaultCRService.GenerateDefaultCRFile(opts.Out, defaultCRFilePath); err != nil {
			return err
		}
	}

	return nil
}
