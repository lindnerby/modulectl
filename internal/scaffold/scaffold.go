package scaffold

import (
	"path"

	"github.com/kyma-project/modulectl/tools/io"
)

type ModuleConfigService interface {
	PreventOverwrite(directory, moduleConfigFileName string, overwrite bool) error
}

type ManifestService interface {
	GenerateManifestFile(out io.Out, manifestFilePath string) error
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
	if err := s.manifestService.GenerateManifestFile(opts.Out, manifestFilePath); err != nil {
		return err
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
