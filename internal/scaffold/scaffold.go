package scaffold

import (
	"fmt"
	"path"

	"github.com/kyma-project/modulectl/internal/scaffold/common/types"
	"github.com/kyma-project/modulectl/internal/scaffold/contentprovider"
	"github.com/kyma-project/modulectl/tools/io"
)

type ModuleConfigService interface {
	PreventOverwrite(directory, moduleConfigFileName string, overwrite bool) error
}

type ManifestService interface {
	GenerateManifestFile(out io.Out, path string) error
}

type DefaultCRService interface {
	GenerateDefaultCRFile(out io.Out, path string) error
}

type FileGeneratorService interface {
	GenerateFile(out io.Out, path string, args types.KeyValueArgs) error
}

type ScaffoldService struct {
	moduleConfigService   ModuleConfigService
	manifestService       ManifestService
	defaultCRService      DefaultCRService
	securityConfigService FileGeneratorService
}

func NewScaffoldService(moduleConfigService ModuleConfigService,
	manifestService ManifestService,
	defaultCRService DefaultCRService,
	securityConfigService FileGeneratorService) *ScaffoldService {
	return &ScaffoldService{
		moduleConfigService:   moduleConfigService,
		manifestService:       manifestService,
		defaultCRService:      defaultCRService,
		securityConfigService: securityConfigService,
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

	securityConfigFilePath := ""
	if opts.securityConfigFileNameConfigured() {
		securityConfigFilePath = path.Join(opts.Directory, opts.SecurityConfigFileName)
		if err := s.securityConfigService.GenerateFile(
			opts.Out,
			securityConfigFilePath,
			types.KeyValueArgs{contentprovider.ArgModuleName: opts.ModuleName}); err != nil {
			return fmt.Errorf("%w %s: %w", ErrGenertingFile, opts.SecurityConfigFileName, err)
		}
	}

	return nil
}
