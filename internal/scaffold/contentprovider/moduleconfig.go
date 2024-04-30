package contentprovider

import (
	"fmt"

	"github.com/kyma-project/modulectl/internal/scaffold/common/types"
)

type ModuleConfigContentProvider struct {
	yamlConverter ObjectToYAMLConverter
}

func NewModuleConfigContentProvider(yamlConverter ObjectToYAMLConverter) *ModuleConfigContentProvider {
	return &ModuleConfigContentProvider{
		yamlConverter: yamlConverter,
	}
}

func (s *ModuleConfigContentProvider) GetDefaultContent(args types.KeyValueArgs) (string, error) {
	if err := s.validateArgs(args); err != nil {
		return "", err
	}

	return "not implemented yet", nil
	// return s.yamlConverter.ConvertToYaml("NotImplemented"), nil
}

func (s *ModuleConfigContentProvider) validateArgs(args types.KeyValueArgs) error {
	if args == nil {
		return fmt.Errorf("%w: args must not be nil", ErrInvalidArg)
	}

	if value, ok := args[ArgModuleName]; !ok {
		return fmt.Errorf("%w: %s", ErrMissingArg, ArgModuleName)
	} else if value == "" {
		return fmt.Errorf("%w: %s must not be empty", ErrInvalidArg, ArgModuleName)
	}

	if value, ok := args[ArgModuleVersion]; !ok {
		return fmt.Errorf("%w: %s", ErrMissingArg, ArgModuleVersion)
	} else if value == "" {
		return fmt.Errorf("%w: %s must not be empty", ErrInvalidArg, ArgModuleVersion)
	}

	if value, ok := args[ArgModuleChannel]; !ok {
		return fmt.Errorf("%w: %s", ErrMissingArg, ArgModuleChannel)
	} else if value == "" {
		return fmt.Errorf("%w: %s must not be empty", ErrInvalidArg, ArgModuleChannel)
	}

	return nil
}
