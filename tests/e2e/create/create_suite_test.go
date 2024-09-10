//go:build e2e

package create_test

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
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
	path                          string
	version                       string
	moduleConfigFile              string
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

	if cmd.version != "" {
		args = append(args, "--version="+cmd.version)
	}
	if cmd.moduleArchiveVersionOverwrite {
		args = append(args, "--module-archive-version-overwrite")
	}

	command = exec.Command("kyma", args...)
	cmdOut, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("create command failed with output: %s and error: %w", cmdOut, err)
	}
	return nil
}

func filesIn(dir string) []string {
	fi, err := os.Stat(dir)
	Expect(err).ToNot(HaveOccurred())
	Expect(fi.IsDir()).To(BeTrueBecause("The provided path should be a directory: %s", dir))

	dirFs := os.DirFS(dir)
	entries, err := fs.ReadDir(dirFs, ".")
	Expect(err).ToNot(HaveOccurred())

	var res []string
	for _, ent := range entries {
		if ent.Type().IsRegular() {
			res = append(res, ent.Name())
		}
	}

	return res
}

func resolveWorkingDirectory() string {
	createDir := os.Getenv("CREATE_DIR")
	if len(createDir) > 0 {
		return createDir
	}

	createDir, err := os.MkdirTemp("", "create_test")
	if err != nil {
		Fail(err.Error())
	}
	return createDir
}

func cleanupWorkingDirectory(path string) {
	if len(os.Getenv("CREATE_DIR")) == 0 {
		_ = os.RemoveAll(path)
	}
}
