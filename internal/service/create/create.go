package create

type ModuleConfigService interface{}

type Service struct {
	moduleConfigService ModuleConfigService
}

func NewService(moduleConfigService ModuleConfigService) (*Service, error) {
	return &Service{
		moduleConfigService: moduleConfigService,
	}, nil
}

func (s *Service) Run(opts Options) error {
	if err := opts.Validate(); err != nil {
		return err
	}
	return nil
}
