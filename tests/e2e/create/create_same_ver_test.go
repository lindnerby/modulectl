//go:build e2e

package create_test

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errCreateModuleFailedWithSameVersion = errors.New(
	"failed to create module with same version exists message")

const (
	ociRepositoryEnvVar         = "OCI_REPOSITORY_URL"
	moduleTemplateVersionEnvVar = "MODULE_TEMPLATE_VERSION"
)

func Test_SameVersion_ModuleCreation(t *testing.T) {
	path := "../../../template-operator"
	configFilePath := "../../../template-operator/module-config.yaml"
	secScannerConfigFile := "../../../template-operator/sec-scanners-config.yaml"
	changedSecScannerConfigFile := "../../../template-operator/sec-scanners-config-changed.yaml"
	version := os.Getenv(moduleTemplateVersionEnvVar)
	registry := os.Getenv(ociRepositoryEnvVar)

	t.Run("Create same version module with module-archive-version-overwrite flag", func(t *testing.T) {
		err := createModuleCommand(true, path, registry, configFilePath, version, secScannerConfigFile)
		assert.NoError(t, err)
	})

	t.Run("Create same version module and same content without module-archive-version-overwrite flag",
		func(t *testing.T) {
			err := createModuleCommand(false, path, registry, configFilePath, version, secScannerConfigFile)
			assert.NoError(t, err)
		})

	t.Run("Create same version module, but different content without module-archive-version-overwrite flag",
		func(t *testing.T) {
			err := createModuleCommand(false, path, registry, configFilePath, version, changedSecScannerConfigFile)
			assert.Equal(t, errCreateModuleFailedWithSameVersion, err)
		})
}

func createModuleCommand(versionOverwrite bool,
	path, registry, configFilePath, version, secScannerConfig string,
) error {
	var createModuleCmd *exec.Cmd
	if versionOverwrite {
		createModuleCmd = exec.Command("kyma", "alpha", "create", "module",
			"--path", path, "--registry", registry, "--insecure", "--module-config-file", configFilePath,
			"--version", version, "--module-archive-version-overwrite", "--sec-scanners-config", secScannerConfig)
	} else {
		createModuleCmd = exec.Command("kyma", "alpha", "create", "module",
			"--path", path, "--registry", registry, "--insecure", "--module-config-file", configFilePath,
			"--version", version, "--sec-scanners-config", secScannerConfig)
	}
	createOut, err := createModuleCmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(createOut),
			fmt.Sprintf("version %s already exists with different content", version)) {
			return errCreateModuleFailedWithSameVersion
		}
		return fmt.Errorf("create module command failed: %w", err)
	}
	return nil
}
