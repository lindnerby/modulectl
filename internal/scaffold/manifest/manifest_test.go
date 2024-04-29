package manifest_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kyma-project/modulectl/internal/scaffold/manifest"
)

func Test_DefaultCRService_GenerateDefaultCRFile_ReturnsError_WhenFileDoesNotExist(t *testing.T) {
	svc := manifest.NewManifestService(&fileExistsErrorStub{})

	result := svc.GenerateManifestFile(&testOut{}, "some-path")

	require.ErrorIs(t, result, manifest.ErrGeneratingManifestFile)
	require.ErrorIs(t, result, errSomeOSError)
}

func Test_DefaultCRService_GenerateDefaultCRFile_Returns_WhenFileDoesAlreadyExist(t *testing.T) {
	out := &testOut{}
	svc := manifest.NewManifestService(&fileExistsStub{})

	result := svc.GenerateManifestFile(out, "some-path")

	require.NoError(t, result)
	require.Len(t, out.sink, 1)
	assert.Contains(t, out.sink[0], "The manifest file already exists, reusing:")
}

func Test_DefaultCRService_GenerateDefaultCRFile_ReturnsError_WhenErrorWritingFile(t *testing.T) {
	out := &testOut{}
	svc := manifest.NewManifestService(&fileWriteErrorStub{})

	result := svc.GenerateManifestFile(out, "some-path")

	require.ErrorIs(t, result, manifest.ErrGeneratingManifestFile)
	require.ErrorIs(t, result, manifest.ErrWritingManifestFile)
	require.Len(t, out.sink, 0)
}

func Test_DefaultCRService_GenerateDefaultCRFile_Returns_WhenFileIsGenerated(t *testing.T) {
	out := &testOut{}
	svc := manifest.NewManifestService(&fileDoesNotExistStub{})

	result := svc.GenerateManifestFile(out, "some-path")

	require.NoError(t, result)
	require.Len(t, out.sink, 1)
	assert.Contains(t, out.sink[0], "Generated a blank manifest file:")
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

type fileExistsErrorStub struct{}

var errSomeOSError = errors.New("some os error")

func (*fileExistsErrorStub) WriteFile(_, _ string) error {
	return nil
}

func (*fileExistsErrorStub) FileExists(_ string) (bool, error) {
	return false, errSomeOSError
}

type fileWriteErrorStub struct{}

func (*fileWriteErrorStub) WriteFile(_, _ string) error {
	return errSomeOSError
}

func (*fileWriteErrorStub) FileExists(_ string) (bool, error) {
	return false, nil
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
