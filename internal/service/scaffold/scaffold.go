package scaffold

import (
	"fmt"
	"path"

	"github.com/kyma-project/modulectl/internal/scaffold/common/errors"
	"github.com/kyma-project/modulectl/internal/scaffold/common/types"
	"github.com/kyma-project/modulectl/internal/service/contentprovider"
	"github.com/kyma-project/modulectl/tools/io"
)

type ModuleConfigService interface {
	FileGeneratorService
	PreventOverwrite(directory, moduleConfigFileName string, overwrite bool) error
}

type FileGeneratorService interface {
	GenerateFile(out io.Out, path string, args types.KeyValueArgs) error
}

type Service struct {
	moduleConfigService   ModuleConfigService
	manifestService       FileGeneratorService
	defaultCRService      FileGeneratorService
	securityConfigService FileGeneratorService
}

func NewService(moduleConfigService ModuleConfigService,
	manifestService FileGeneratorService,
	defaultCRService FileGeneratorService,
	securityConfigService FileGeneratorService) (*Service, error) {
	if moduleConfigService == nil {
		return nil, fmt.Errorf("%w: moduleConfigService must not be nil", errors.ErrInvalidArg)
	}

	if manifestService == nil {
		return nil, fmt.Errorf("%w: manifestService must not be nil", errors.ErrInvalidArg)
	}

	if defaultCRService == nil {
		return nil, fmt.Errorf("%w: defaultCRService must not be nil", errors.ErrInvalidArg)
	}

	if securityConfigService == nil {
		return nil, fmt.Errorf("%w: securityConfigService must not be nil", errors.ErrInvalidArg)
	}

	return &Service{
		moduleConfigService:   moduleConfigService,
		manifestService:       manifestService,
		defaultCRService:      defaultCRService,
		securityConfigService: securityConfigService,
	}, nil
}

func (s *Service) CreateScaffold(opts Options) error {
	if err := opts.validate(); err != nil {
		return err
	}

	if err := s.moduleConfigService.PreventOverwrite(opts.Directory, opts.ModuleConfigFileName, opts.ModuleConfigFileOverwrite); err != nil {
		return err
	}

	manifestFilePath := path.Join(opts.Directory, opts.ManifestFileName)
	if err := s.manifestService.GenerateFile(opts.Out, manifestFilePath, nil); err != nil {
		return fmt.Errorf("%w %s: %w", ErrGeneratingFile, opts.ManifestFileName, err)
	}

	defaultCRFilePath := ""
	if opts.defaultCRFileNameConfigured() {
		defaultCRFilePath = path.Join(opts.Directory, opts.DefaultCRFileName)
		if err := s.defaultCRService.GenerateFile(opts.Out, defaultCRFilePath, nil); err != nil {
			return fmt.Errorf("%w %s: %w", ErrGeneratingFile, opts.DefaultCRFileName, err)
		}
	}

	securityConfigFilePath := ""
	if opts.securityConfigFileNameConfigured() {
		securityConfigFilePath = path.Join(opts.Directory, opts.SecurityConfigFileName)
		if err := s.securityConfigService.GenerateFile(
			opts.Out,
			securityConfigFilePath,
			types.KeyValueArgs{contentprovider.ArgModuleName: opts.ModuleName}); err != nil {
			return fmt.Errorf("%w %s: %w", ErrGeneratingFile, opts.SecurityConfigFileName, err)
		}
	}

	moduleConfigFilePath := path.Join(opts.Directory, opts.ModuleConfigFileName)
	if err := s.moduleConfigService.GenerateFile(
		opts.Out,
		moduleConfigFilePath,
		types.KeyValueArgs{
			contentprovider.ArgModuleName:         opts.ModuleName,
			contentprovider.ArgModuleVersion:      opts.ModuleVersion,
			contentprovider.ArgModuleChannel:      opts.ModuleChannel,
			contentprovider.ArgManifestFile:       opts.ManifestFileName,
			contentprovider.ArgDefaultCRFile:      defaultCRFilePath,
			contentprovider.ArgSecurityConfigFile: securityConfigFilePath,
		}); err != nil {
		return fmt.Errorf("%w %s: %w", ErrGeneratingFile, opts.ModuleConfigFileName, err)
	}

	return nil
}
