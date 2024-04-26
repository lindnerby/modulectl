package scaffold_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kyma-project/modulectl/internal/scaffold"
	"github.com/kyma-project/modulectl/internal/scaffold/moduleconfig"
	"github.com/kyma-project/modulectl/internal/testutils/builder"
	"github.com/kyma-project/modulectl/internal/testutils/stub"
)

func Test_RunScaffold_ReturnsError_WhenOutIsNil(t *testing.T) {
	scaffoldService := getScaffoldService()
	opts := builder.NewScaffoldOptionsBuilder().WithOut(nil).Build()

	result := scaffoldService.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.Out")
}

func Test_RunScaffold_ReturnsError_WhenDirectoryIsEmpty(t *testing.T) {
	scaffoldService := getScaffoldService()
	opts := builder.NewScaffoldOptionsBuilder().WithDirectory("").Build()

	result := scaffoldService.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.Directory")
}

func Test_RunScaffold_ReturnsError_WhenModuleConfigFileIsEmpty(t *testing.T) {
	scaffoldService := getScaffoldService()
	opts := builder.NewScaffoldOptionsBuilder().WithModuleConfigFileName("").Build()

	result := scaffoldService.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ModuleConfigFileName")
}

func Test_RunScaffold_ReturnsError_WhenFilesystemReturnsError(t *testing.T) {
	scaffoldService := scaffold.NewScaffoldService(
		moduleconfig.NewModuleConfigService(
			&stub.ErrorStub{},
		),
	)
	opts := builder.NewScaffoldOptionsBuilder().Build()

	result := scaffoldService.CreateScaffold(opts)

	require.ErrorIs(t, result, stub.ErrSomeOSError)
}

func Test_RunScaffold_ReturnsError_WhenModuleConfigFileExistsButOverwriteIsNotSet(t *testing.T) {
	scaffoldService := scaffold.NewScaffoldService(
		moduleconfig.NewModuleConfigService(
			&stub.FileExistsStub{},
		),
	)
	opts := builder.NewScaffoldOptionsBuilder().WithModuleConfigFileOverwrite(false).Build()

	result := scaffoldService.CreateScaffold(opts)

	require.ErrorIs(t, result, moduleconfig.ErrFileExists)
}

func Test_RunScaffold_Succeeds_WhenModuleConfigFileExistsAndOverwriteIsSet(t *testing.T) {
	scaffoldService := scaffold.NewScaffoldService(
		moduleconfig.NewModuleConfigService(
			&stub.FileExistsStub{},
		),
	)
	opts := builder.NewScaffoldOptionsBuilder().WithModuleConfigFileOverwrite(true).Build()

	result := scaffoldService.CreateScaffold(opts)

	require.NoError(t, result)
}

func getScaffoldService() *scaffold.ScaffoldService {
	return scaffold.NewScaffoldService(
		moduleconfig.NewModuleConfigService(
			&stub.FileDoesNotExistStub{},
		),
	)
}
