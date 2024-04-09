package scaffold

import (
	"fmt"
	"github.com/kyma-project/modulectl/tools/common"
	"os"
	"path"
)

func ValidateFlags() error {
	if ModuleName == "" {
		return errModuleNameEmpty
	}

	if ModuleVersion == "" {
		return errModuleVersionEmpty
	}

	if ModuleChannel == "" {
		return errModuleChannelEmpty
	}

	err := common.ValidateDirectory(Directory)
	if err != nil {
		return err
	}

	if ModuleConfigFile == "" {
		return errModuleConfigEmpty
	}

	if ManifestFile == "" {
		return errManifestFileEmpty
	}

	return nil
}

// ***********************
// **** ModuleConfig *****
// ***********************

func ModuleConfigFilePath() string {
	return path.Join(Directory, ModuleConfigFile)
}

func ModuleConfigFileExists() (bool, error) {
	return common.FileExists(ModuleConfigFilePath())
}

func GenerateModuleConfigFile() error {
	cfg := Config{
		Name:          ModuleName,
		Version:       ModuleVersion,
		Channel:       ModuleChannel,
		ManifestPath:  ManifestFile,
		Security:      SecurityConfigFile,
		DefaultCRPath: DefaultCRFile,
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	return common.GenerateYamlFileFromObject(cfg, ModuleConfigFilePath())
}

// ***********************
// ****** Manifest *******
// ***********************

func ManifestFilePath() string {
	return path.Join(Directory, ManifestFile)
}

func ManifestFileExists() (bool, error) {
	return common.FileExists(ManifestFilePath())
}

func GenerateManifest() error {
	blankContents := `# This file holds the Manifest of your module, encompassing all resources installed in the cluster once the module is activated.
# It should include the Custom Resource Definition for your module's default CustomResource, if it exists.

`
	filePath := ManifestFilePath()
	err := os.WriteFile(filePath, []byte(blankContents), 0600)
	if err != nil {
		return fmt.Errorf("error while saving %s: %w", filePath, err)
	}

	return nil
}

// ***********************
// ****** DefaultCR ******
// ***********************

func defaultCRFileConfigured() bool {
	return DefaultCRFile != ""
}

func DefaultCRFilePath() string {
	return path.Join(Directory, DefaultCRFile)
}

func DefaultCRFileExists() (bool, error) {
	return common.FileExists(DefaultCRFilePath())
}

func GenerateDefaultCRFile() error {
	blankContents := `# This is the file that contains the defaultCR for your module, which is the Custom Resource that will be created upon module enablement.
# Make sure this file contains *ONLY* the Custom Resource (not the Custom Resource Definition, which should be a part of your module manifest)

`
	filePath := DefaultCRFilePath()
	err := os.WriteFile(filePath, []byte(blankContents), 0600)
	if err != nil {
		return fmt.Errorf("error while saving %s: %w", filePath, err)
	}

	return nil
}

// ***********************
// ****** SecConfig ******
// ***********************

type SecurityScanCfg struct {
	ModuleName  string            `json:"module-name" yaml:"module-name" comment:"string, name of your module"`
	Protecode   []string          `json:"protecode" yaml:"protecode" comment:"list, includes the images which must be scanned by the Protecode scanner (aka. Black Duck Binary Analysis)"`
	WhiteSource WhiteSourceSecCfg `json:"whitesource" yaml:"whitesource" comment:"whitesource (aka. Mend) security scanner specific configuration"`
	DevBranch   string            `json:"dev-branch" yaml:"dev-branch" comment:"string, name of the development branch"`
	RcTag       string            `json:"rc-tag" yaml:"rc-tag" comment:"string, release candidate tag"`
}

type WhiteSourceSecCfg struct {
	Language    string   `json:"language" yaml:"language" comment:"string, indicating the programming language the scanner has to analyze"`
	SubProjects string   `json:"subprojects" yaml:"subprojects" comment:"string, specifying any subprojects"`
	Exclude     []string `json:"exclude" yaml:"exclude" comment:"list, directories within the repository which should not be scanned"`
}

func securityConfigFileConfigured() bool {
	return SecurityConfigFile != ""
}

func SecurityConfigFilePath() string {
	return path.Join(Directory, SecurityConfigFile)
}

func SecurityConfigFileExists() (bool, error) {
	return common.FileExists(SecurityConfigFilePath())
}

func GenerateSecurityConfigFile() error {
	cfg := SecurityScanCfg{
		ModuleName: ModuleName,
		Protecode: []string{"europe-docker.pkg.dev/kyma-project/prod/myimage:1.2.3",
			"europe-docker.pkg.dev/kyma-project/prod/external/ghcr.io/mymodule/anotherimage:4.5.6"},
		WhiteSource: WhiteSourceSecCfg{
			Exclude: []string{"**/test/**", "**/*_test.go"},
		},
	}
	err := common.GenerateYamlFileFromObject(cfg, SecurityConfigFilePath())
	return err
}
