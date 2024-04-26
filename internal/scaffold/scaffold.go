package scaffold

import (
	"github.com/kyma-project/modulectl/internal/scaffold/moduleconfig"
)

type ScaffoldService struct {
	moduleConfigService *moduleconfig.ModuleConfigService
}

func NewScaffoldService(moduleConfigService *moduleconfig.ModuleConfigService) *ScaffoldService {
	return &ScaffoldService{
		moduleConfigService: moduleConfigService,
	}
}

func (s *ScaffoldService) CreateScaffold(opts Options) error {
	if err := opts.validate(); err != nil {
		return err
	}

	if err := s.moduleConfigService.PreventOverwrite(opts.Directory, opts.ModuleConfigFileName, opts.ModuleConfigFileOverwrite); err != nil {
		return err
	}

	return nil
}
