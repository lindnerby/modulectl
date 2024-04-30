package filegenerator_test

import (
	"errors"
	"testing"

	"github.com/kyma-project/modulectl/internal/scaffold/common/types"
	"github.com/kyma-project/modulectl/internal/scaffold/filegenerator"
	"github.com/kyma-project/modulectl/tools/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ReuseFileGeneratorService_GenerateFile_ReturnsError_WhenFileExistenceCheckReturnsError(t *testing.T) {
	svc := filegenerator.NewReuseFileGeneratorService("test-kind", &fileExistsErrorStub{}, &fileGeneratorErrorStub{})

	result := svc.GenerateFile(&rfgTestOut{}, "some-path", nil)

	require.ErrorIs(t, result, filegenerator.ErrCheckingFileExistence)
	require.ErrorIs(t, result, errSomeFileExistsOSError)
}

func Test_ReuseFileGeneratorService_GenerateFile_Succeeds_WhenFileDoesAlreadyExist(t *testing.T) {
	out := &rfgTestOut{}
	svc := filegenerator.NewReuseFileGeneratorService("test-kind", &fileExistsStub{}, &fileGeneratorErrorStub{})

	result := svc.GenerateFile(out, "some-path", nil)

	require.NoError(t, result)
	require.Len(t, out.sink, 1)
	assert.Contains(t, out.sink[0], "The test-kind file already exists, reusing:")
}

func Test_ReuseFileGeneratorService_GenerateFile_ReturnsError_WhenFileGenerationReturnsError(t *testing.T) {
	svc := filegenerator.NewReuseFileGeneratorService("test-kind", &fileDoesNotExistStub{}, &fileGeneratorErrorStub{})

	result := svc.GenerateFile(&rfgTestOut{}, "some-path", nil)

	require.ErrorIs(t, result, errSomeFileGeneratorError)
}

func Test_ReuseFileGeneratorService_GenerateFile_Succeeds_WhenFileIsGenerated(t *testing.T) {
	svc := filegenerator.NewReuseFileGeneratorService("test-kind", &fileExistsStub{}, &fileGeneratorStub{})

	result := svc.GenerateFile(&rfgTestOut{}, "some-path", nil)

	require.NoError(t, result)
}

// ***************
// Test Stubs
// ***************

type rfgTestOut struct {
	sink []string
}

func (o *rfgTestOut) Write(msg string) {
	o.sink = append(o.sink, msg)
}

type fileGeneratorStub struct{}

func (*fileGeneratorStub) GenerateFile(out io.Out, path string, args types.KeyValueArgs) error {
	return nil
}

type fileGeneratorErrorStub struct{}

var errSomeFileGeneratorError = errors.New("some file generator error")

func (*fileGeneratorErrorStub) GenerateFile(out io.Out, path string, args types.KeyValueArgs) error {
	return errSomeFileGeneratorError
}

type fileExistsErrorStub struct{}

var errSomeFileExistsOSError = errors.New("some os error")

func (*fileExistsErrorStub) FileExists(_ string) (bool, error) {
	return false, errSomeFileExistsOSError
}

type fileExistsStub struct{}

func (*fileExistsStub) FileExists(_ string) (bool, error) {
	return true, nil
}

type fileDoesNotExistStub struct{}

func (*fileDoesNotExistStub) FileExists(_ string) (bool, error) {
	return false, nil
}
