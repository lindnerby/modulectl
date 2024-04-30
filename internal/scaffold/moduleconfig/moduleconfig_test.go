package moduleconfig_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kyma-project/modulectl/internal/scaffold/common/types"
	"github.com/kyma-project/modulectl/internal/scaffold/moduleconfig"
	"github.com/kyma-project/modulectl/tools/io"
)

const (
	directory        = "./bin/dir"
	moduleConfigFile = "config.yaml"
)

func Test_PreventOverwrite_ReturnsError_WhenFilesystemReturnsError(t *testing.T) {
	svc := moduleconfig.NewModuleConfigService(
		&errorStub{},
		&fileGeneratorErrorStub{},
	)

	result := svc.PreventOverwrite(directory, moduleConfigFile, true)

	require.ErrorIs(t, result, errSomeOSError)
}

func Test_PreventOverwrite_ReturnsError_WhenFileExistsAndNoOverwrite(t *testing.T) {
	svc := moduleconfig.NewModuleConfigService(
		&fileExistsStub{},
		&fileGeneratorErrorStub{},
	)

	result := svc.PreventOverwrite(directory, moduleConfigFile, false)

	require.ErrorIs(t, result, moduleconfig.ErrFileExists)
}

func Test_PreventOverwrite_ReturnsNil_WhenFileExistsAndOverwrite(t *testing.T) {
	svc := moduleconfig.NewModuleConfigService(
		&fileExistsStub{},
		&fileGeneratorErrorStub{},
	)

	result := svc.PreventOverwrite(directory, moduleConfigFile, true)

	require.NoError(t, result)
}

func Test_PreventOverwrite_ReturnsNil_WhenFileDoesNotExist(t *testing.T) {
	svc := moduleconfig.NewModuleConfigService(
		&fileDoesNotExistStub{},
		&fileGeneratorErrorStub{},
	)

	result := svc.PreventOverwrite(directory, moduleConfigFile, true)

	require.NoError(t, result)
}

// ***************
// Test Stubs
// ***************

type fileExistsStub struct{}

func (*fileExistsStub) FileExists(_ string) (bool, error) {
	return true, nil
}

type fileDoesNotExistStub struct{}

func (*fileDoesNotExistStub) FileExists(_ string) (bool, error) {
	return false, nil
}

var errSomeOSError = errors.New("some OS error")

type errorStub struct{}

func (*errorStub) FileExists(_ string) (bool, error) {
	return false, errSomeOSError
}

type fileGeneratorErrorStub struct{}

var errSomeFileGeneratorError = errors.New("some file generator error")

func (*fileGeneratorErrorStub) GenerateFile(out io.Out, path string, args types.KeyValueArgs) error {
	return errSomeFileGeneratorError
}
