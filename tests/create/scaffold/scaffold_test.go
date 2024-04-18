package scaffold

import (
	"bytes"
	"github.com/kyma-project/modulectl/cmd/modulectl/create/scaffold"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

// ***********************
// ***** TEST CASES ******
// ***********************
func (s *ScaffoldCommandLayerSuite) TestNoArgsAllowed() {
	s.scaffoldCmd.SetArgs([]string{"arg1", "arg2", "arg3"})
	err := s.scaffoldCmd.Execute()
	assert.ErrorContains(s.T(), err, "accepts 0 arg(s), received 3")
}

func (s *ScaffoldCommandLayerSuite) TestHelpFlag() {
	s.scaffoldCmd.SetArgs([]string{"--help"})
	err := s.scaffoldCmd.Execute()
	assert.NoError(s.T(), err)

	randomPartOfLongHelpMessage := `The command is designed to streamline the module creation process in Kyma, making it easier and more 
efficient for developers to get started with new modules. It supports customization through various flags, 
allowing for a tailored scaffolding experience according to the specific needs of the module being created.
`

	assert.Contains(s.T(), s.localStdout.String(), randomPartOfLongHelpMessage)

	assert.Contains(s.T(), s.localStdout.String(), "Usage:\n  scaffold")
	assert.Contains(s.T(), s.localStdout.String(), "--module-name=NAME")
	assert.Contains(s.T(), s.localStdout.String(), "--module-version=VERSION")
	assert.Contains(s.T(), s.localStdout.String(), "--module-channel=CHANNEL")
}

func (s *ScaffoldCommandLayerSuite) TestDefaultFlagsCorrectlySet() {
	s.scaffoldCmd.SetArgs([]string{})
	err := s.scaffoldCmd.Execute()
	assert.NoError(s.T(), err)

	for _, test := range getFlagsWithDefaults() {
		assert.Equal(s.T(), test.flagDefault, s.scaffoldCmd.Flags().Lookup(test.flagName).Value.String())
		assert.False(s.T(), s.scaffoldCmd.Flags().Lookup(test.flagName).Changed)
	}
}

func (s *ScaffoldCommandLayerSuite) TestDefaultFlagsAsInput() {
	s.scaffoldCmd.SetArgs([]string{"--module-name=kyma-project.io/module/mymodule", "--module-version=0.0.1", "--module-channel=regular"})
	err := s.scaffoldCmd.Execute()
	assert.NoError(s.T(), err)

	for _, test := range getFlagsWithDefaults() {
		assert.Equal(s.T(), test.flagDefault, s.scaffoldCmd.Flags().Lookup(test.flagName).Value.String())
		assert.True(s.T(), s.scaffoldCmd.Flags().Lookup(test.flagName).Changed)
	}
}

func (s *ScaffoldCommandLayerSuite) TestNoFlagsSet() {
	s.scaffoldCmd.SetArgs([]string{})
	err := s.scaffoldCmd.Execute()
	assert.NoError(s.T(), err)

	outputMessages := []string{
		"Validating",
		"Generating Manifest file",
		"Generated a blank Manifest file: manifest.yaml",
		"Generating module config file",
		"Generated module config file: scaffold-module-config.yaml",
	}

	for _, outputMessage := range outputMessages {
		assert.Contains(s.T(), s.localStdout.String(), outputMessage)
	}
}

func (s *ScaffoldCommandLayerSuite) TestOverwriteFlag() {
	s.cleanGeneratedFiles()

	s.scaffoldCmd.SetArgs([]string{})
	err := s.scaffoldCmd.Execute()
	assert.NoError(s.T(), err)

	assert.FileExists(s.T(), "manifest.yaml")
	assert.FileExists(s.T(), "scaffold-module-config.yaml")

	err = s.scaffoldCmd.Execute()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "module config file already exists. use --overwrite flag to overwrite it")
}

//	More potential test cases
//func (s *ScaffoldCommandLayerSuite) TestAllFlagsSet() {}

//func (s *ScaffoldCommandLayerSuite) TestValidInputs() {}
//
//func (s *ScaffoldCommandLayerSuite) TestInvalidInputs() {
//	//	eg: missing flag values, etc.
//	//	eg: case sensitivity of the flags
//}

// ***********************
// ******** SETUP ********
// ***********************
type ScaffoldCommandLayerSuite struct {
	suite.Suite
	scaffoldCmd *cobra.Command
	localStdout bytes.Buffer
}

func (s *ScaffoldCommandLayerSuite) SetupTest() {
	s.scaffoldCmd = scaffold.NewCmd()
	s.localStdout = bytes.Buffer{}

	s.scaffoldCmd.SetOut(&s.localStdout)

	s.cleanGeneratedFiles()
}

func (s *ScaffoldCommandLayerSuite) cleanGeneratedFiles() {
	err := os.Remove("./manifest.yaml")
	if err != nil {
		return
	}
	err = os.Remove("./scaffold-module-config.yaml")
	if err != nil {
		return
	}
}

func (s *ScaffoldCommandLayerSuite) AfterTest() {
	s.cleanGeneratedFiles()
}

func TestScaffoldCommandSuite(t *testing.T) {
	suite.Run(t, new(ScaffoldCommandLayerSuite))
}

type FlagDefaults struct {
	flagName    string
	flagDefault string
}

func getFlagsWithDefaults() []FlagDefaults {
	return []FlagDefaults{

		{
			flagName:    "module-name",
			flagDefault: "kyma-project.io/module/mymodule",
		},
		{
			flagName:    "module-version",
			flagDefault: "0.0.1",
		},
		{
			flagName:    "module-channel",
			flagDefault: "regular",
		},
	}
}
