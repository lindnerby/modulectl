package scaffold_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/kyma-project/modulectl/cmd/modulectl/create/scaffold"
)

func Test_ScaffoldFlagsDefaults(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{name: scaffold.DirectoryFlagName, value: scaffold.DirectoryFlagDefault, expected: "./"},
		{name: scaffold.ModuleConfigFileFlagName, value: scaffold.ModuleConfigFileFlagDefault, expected: "scaffold-module-config.yaml"},
		{name: scaffold.ModuleConfigFileOverwriteFlagName, value: strconv.FormatBool(scaffold.ModuleConfigFileOverwriteFlagDefault), expected: "false"},
		{name: scaffold.ManifestFileFlagName, value: scaffold.ManifestFileFlagDefault, expected: "manifest.yaml"},
		{name: scaffold.DefaultCRFlagName, value: scaffold.DefaultCRFlagDefault, expected: "default-cr.yaml"},
	}

	for _, testcase := range tests {
		testcase := testcase
		testName := fmt.Sprintf("TestFlagHasCorrectDefault_%s", testcase.name)
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			if testcase.value != testcase.expected {
				t.Errorf("Flag '%s' has different default: expected = '%s', got = '%s'",
					testcase.name, testcase.expected, testcase.value)
			}
		})
	}
}
