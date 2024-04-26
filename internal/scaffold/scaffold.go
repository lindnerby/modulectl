package scaffold

type ModuleConfigService interface {
	PreventOverwrite(directory, moduleConfigFileName string, overwrite bool) error
}

type ScaffoldService struct {
	moduleConfigService ModuleConfigService
}

func NewScaffoldService(moduleConfigService ModuleConfigService) *ScaffoldService {
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
