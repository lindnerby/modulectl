package contentprovider

import (
	"fmt"

	"github.com/kyma-project/modulectl/internal/scaffold/common/types"
)

const (
	ArgModuleName = "moduleName"
)

type ObjectToYAMLConverter interface {
	ConvertToYaml(obj interface{}) string
}

type SecurityConfigContentProvider struct {
	yamlConverter ObjectToYAMLConverter
}

func NewSecurityConfigContentProvider(yamlConverter ObjectToYAMLConverter) *SecurityConfigContentProvider {
	return &SecurityConfigContentProvider{
		yamlConverter: yamlConverter,
	}
}

func (s *SecurityConfigContentProvider) GetDefaultContent(args types.KeyValueArgs) (string, error) {
	if err := validateArgs(args); err != nil {
		return "", err
	}

	return s.yamlConverter.ConvertToYaml(getSecurityConfig(args[ArgModuleName])), nil
}

func validateArgs(args types.KeyValueArgs) error {
	if args == nil {
		return fmt.Errorf("%w: args must not be nil", ErrInvalidArg)
	}

	value, ok := args[ArgModuleName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrMissingArg, ArgModuleName)
	}

	if value == "" {
		return fmt.Errorf("%w: %s must not be empty", ErrInvalidArg, ArgModuleName)
	}

	return nil
}

func getSecurityConfig(moduleName string) securityScanCfg {
	return securityScanCfg{
		ModuleName: moduleName,
		Protecode: []string{"europe-docker.pkg.dev/kyma-project/prod/myimage:1.2.3",
			"europe-docker.pkg.dev/kyma-project/prod/external/ghcr.io/mymodule/anotherimage:4.5.6"},
		WhiteSource: whiteSourceSecCfg{
			Exclude: []string{"**/test/**", "**/*_test.go"},
		},
	}
}

type securityScanCfg struct {
	ModuleName  string            `json:"module-name" yaml:"module-name" comment:"string, name of your module"`
	Protecode   []string          `json:"protecode" yaml:"protecode" comment:"list, includes the images which must be scanned by the Protecode scanner (aka. Black Duck Binary Analysis)"`
	WhiteSource whiteSourceSecCfg `json:"whitesource" yaml:"whitesource" comment:"whitesource (aka. Mend) security scanner specific configuration"`
	DevBranch   string            `json:"dev-branch" yaml:"dev-branch" comment:"string, name of the development branch"`
	RcTag       string            `json:"rc-tag" yaml:"rc-tag" comment:"string, release candidate tag"`
}

type whiteSourceSecCfg struct {
	Language    string   `json:"language" yaml:"language" comment:"string, indicating the programming language the scanner has to analyze"`
	SubProjects string   `json:"subprojects" yaml:"subprojects" comment:"string, specifying any subprojects"`
	Exclude     []string `json:"exclude" yaml:"exclude" comment:"list, directories within the repository which should not be scanned"`
}
