package manifest_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kyma-project/modulectl/internal/scaffold/manifest"
)

const expectedDefaultContent = `# This file holds the Manifest of your module, encompassing all resources installed in the cluster once the module is activated.
# It should include the Custom Resource Definition for your module's default CustomResource, if it exists.

`

func Test_GetDefaultManifestContent_ReturnsExpectedContent(t *testing.T) {
	svc := manifest.NewManifestService(
		&writeFileErrorStub{},
	)

	result := svc.GetDefaultManifestContent()

	require.Equal(t, expectedDefaultContent, result)
}

func Test_WriteManifestFile_Succeeds(t *testing.T) {
	svc := manifest.NewManifestService(
		&writeFileStub{},
	)

	result := svc.WriteManifestFile("content", "path")

	require.NoError(t, result)
}

func Test_WriteManifestFile_ReturnsError(t *testing.T) {
	svc := manifest.NewManifestService(
		&writeFileErrorStub{},
	)

	result := svc.WriteManifestFile("content", "path")

	require.ErrorIs(t, result, manifest.ErrWritingManifestFile)
	require.ErrorIs(t, result, errSomeOSError)
}

// ***************
// Test Stubs
// ***************

type writeFileErrorStub struct{}

var errSomeOSError = errors.New("some os error")

func (*writeFileErrorStub) WriteFile(_, _ string) error {
	return errSomeOSError
}

type writeFileStub struct{}

func (*writeFileStub) WriteFile(_, _ string) error {
	return nil
}
