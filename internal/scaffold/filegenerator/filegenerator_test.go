package filegenerator_test

import (
	"errors"
	"testing"

	"github.com/kyma-project/modulectl/internal/scaffold/common/types"
	"github.com/kyma-project/modulectl/internal/scaffold/filegenerator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FileGeneratorService_GenerateFile_ReturnsError_WhenFileDoesNotExist(t *testing.T) {
	svc := filegenerator.NewFileGeneratorService("test-kind", &fileExistsErrorStub{}, &contentProviderErrorStub{})

	result := svc.GenerateFile(&testOut{}, "some-path", nil)

	require.ErrorIs(t, result, filegenerator.ErrCheckingFileExistence)
	require.ErrorIs(t, result, errSomeOSError)
}

func Test_FileGeneratorService_GenerateFile_Succeeds_WhenFileDoesAlreadyExist(t *testing.T) {
	out := &testOut{}
	svc := filegenerator.NewFileGeneratorService("test-kind", &fileExistsStub{}, &contentProviderErrorStub{})

	result := svc.GenerateFile(out, "some-path", nil)

	require.NoError(t, result)
	require.Len(t, out.sink, 1)
	assert.Contains(t, out.sink[0], "The test-kind file already exists, reusing:")
}

func Test_FileGeneratorService_GenerateFile_ReturnsError_WhenErrorGettingDefaultContent(t *testing.T) {
	out := &testOut{}
	svc := filegenerator.NewFileGeneratorService("test-kind", &fileWriteErrorStub{}, &contentProviderErrorStub{})

	result := svc.GenerateFile(out, "some-path", nil)

	require.ErrorIs(t, result, filegenerator.ErrGettingDefaultContent)
	require.ErrorIs(t, result, errSomeContentProviderError)
	require.Len(t, out.sink, 0)
}

func Test_FileGeneratorService_GenerateFile_ReturnsError_WhenErrorWritingFile(t *testing.T) {
	out := &testOut{}
	svc := filegenerator.NewFileGeneratorService("test-kind", &fileWriteErrorStub{}, &contentProviderStub{})

	result := svc.GenerateFile(out, "some-path", nil)

	require.ErrorIs(t, result, filegenerator.ErrWritingFile)
	require.Len(t, out.sink, 0)
}

func Test_FileGeneratorService_GenerateFile_Succeeds_WhenFileIsGenerated(t *testing.T) {
	out := &testOut{}
	svc := filegenerator.NewFileGeneratorService("test-kind", &fileDoesNotExistStub{}, &contentProviderStub{})

	result := svc.GenerateFile(out, "some-path", nil)

	require.NoError(t, result)
	require.Len(t, out.sink, 1)
	assert.Contains(t, out.sink[0], "Generated a blank test-kind file:")
}

// ***************
// Test Stubs
// ***************

type testOut struct {
	sink []string
}

func (o *testOut) Write(msg string) {
	o.sink = append(o.sink, msg)
}

type contentProviderStub struct{}

func (*contentProviderStub) GetDefaultContent(args types.KeyValueArgs) (string, error) {
	return "test-content", nil
}

type contentProviderErrorStub struct{}

var errSomeContentProviderError = errors.New("some error")

func (*contentProviderErrorStub) GetDefaultContent(args types.KeyValueArgs) (string, error) {
	return "", errSomeContentProviderError
}

type fileExistsErrorStub struct{}

var errSomeOSError = errors.New("some os error")

func (*fileExistsErrorStub) WriteFile(_, _ string) error {
	return nil
}

func (*fileExistsErrorStub) FileExists(_ string) (bool, error) {
	return false, errSomeOSError
}

type fileExistsStub struct{}

func (*fileExistsStub) WriteFile(_, _ string) error {
	return nil
}

func (*fileExistsStub) FileExists(_ string) (bool, error) {
	return true, nil
}

type fileDoesNotExistStub struct{}

func (*fileDoesNotExistStub) WriteFile(_, _ string) error {
	return nil
}

func (*fileDoesNotExistStub) FileExists(_ string) (bool, error) {
	return false, nil
}

type fileWriteErrorStub struct{}

func (*fileWriteErrorStub) WriteFile(_, _ string) error {
	return errSomeOSError
}

func (*fileWriteErrorStub) FileExists(_ string) (bool, error) {
	return false, nil
}
