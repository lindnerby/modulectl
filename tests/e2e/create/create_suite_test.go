//go:build e2e

package create_test

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func Test_Create(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "'Create' Command Test Suite")
}

// Command wrapper for `modulectl create`

type createCmd struct {
	name                          string
	registry                      string
	path                          string
	output                        string
	version                       string
	moduleConfigFile              string
	secScanConfig                 string
	moduleArchiveVersionOverwrite bool
	insecure                      bool
}

func (cmd *createCmd) execute() error {
	var command *exec.Cmd

	// TODO revert to modulectl only debugging against kyma cli bin for verifying tests
	args := []string{"alpha", "create", "module"}

	if cmd.moduleConfigFile != "" {
		args = append(args, "--module-config-file="+cmd.moduleConfigFile)
	}

	if cmd.path != "" {
		args = append(args, "--path="+cmd.path)
	}

	if cmd.name != "" {
		args = append(args, "--name="+cmd.name)
	}

	if cmd.registry != "" {
		args = append(args, "--registry="+cmd.registry)
	}

	if cmd.secScanConfig != "" {
		args = append(args, "--sec-scanners-config="+cmd.secScanConfig)
	}

	if cmd.output != "" {
		args = append(args, "--output="+cmd.output)
	}

	if cmd.version != "" {
		args = append(args, "--version="+cmd.version)
	}

	if cmd.moduleArchiveVersionOverwrite {
		args = append(args, "--module-archive-version-overwrite")
	}

	if cmd.insecure {
		args = append(args, "--insecure")
	}

	println("Running command: modulectl", strings.Join(args, " "))
	// TODO Remove
	args = append(args, "--non-interactive")

	command = exec.Command("kyma", args...)
	cmdOut, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("create command failed with output: %s and error: %w", cmdOut, err)
	}
	return nil
}
