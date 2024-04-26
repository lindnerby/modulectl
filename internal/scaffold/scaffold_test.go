package scaffold_test

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kyma-project/modulectl/internal/scaffold"
	"github.com/kyma-project/modulectl/internal/scaffold/manifest"
	iotools "github.com/kyma-project/modulectl/tools/io"
)

func Test_RunScaffold_ReturnsError_WhenOutIsNil(t *testing.T) {
	svc := scaffold.NewScaffoldService(&preventOverwriteErrorStub{}, &manifestServiceErrorStub{}, &fileDoesNotExistStub{})
	opts := newScaffoldOptionsBuilder().withOut(nil).build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.Out")
}

func Test_RunScaffold_ReturnsError_WhenDirectoryIsEmpty(t *testing.T) {
	svc := scaffold.NewScaffoldService(&preventOverwriteErrorStub{}, &manifestServiceErrorStub{}, &fileDoesNotExistStub{})
	opts := newScaffoldOptionsBuilder().withDirectory("").build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.Directory")
}

func Test_RunScaffold_ReturnsError_WhenModuleConfigFileIsEmpty(t *testing.T) {
	svc := scaffold.NewScaffoldService(&preventOverwriteErrorStub{}, &manifestServiceErrorStub{}, &fileDoesNotExistStub{})
	opts := newScaffoldOptionsBuilder().withModuleConfigFileName("").build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ModuleConfigFileName")
}

func Test_RunScaffold_ReturnsError_WhenManifestFileIsEmpty(t *testing.T) {
	svc := scaffold.NewScaffoldService(&preventOverwriteErrorStub{}, &manifestServiceErrorStub{}, &fileDoesNotExistStub{})
	opts := newScaffoldOptionsBuilder().withManifestFileName("").build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ManifestFileName")
}

func Test_RunScaffold_ReturnsError_WhenModuleConfigServicePreventOverwriteReturnsError(t *testing.T) {
	svc := scaffold.NewScaffoldService(&preventOverwriteErrorStub{}, &manifestServiceErrorStub{}, &fileDoesNotExistStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.ErrorIs(t, result, errOverwriteError)
}

func Test_RunScaffold_ReturnsError_WhenManifestServiceWriteManifestFileReturnsError(t *testing.T) {
	svc := scaffold.NewScaffoldService(&preventOverwriteStub{}, &manifestServiceErrorStub{}, &fileDoesNotExistStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.ErrorIs(t, result, manifest.ErrWritingManifestFile)
}

func Test_RunScaffold_Succeeds_WhenManifestFileDoesNotExist(t *testing.T) {
	svc := scaffold.NewScaffoldService(&preventOverwriteStub{}, &manifestServiceStub{}, &fileDoesNotExistStub{})

	testOut := testOut{sink: []string{}}
	result := svc.CreateScaffold(newScaffoldOptionsBuilder().withOut(&testOut).build())

	require.NoError(t, result)
	require.Len(t, testOut.sink, 1)
	assert.Contains(t, testOut.sink[0], "Generated a blank Manifest file:")
}

func Test_RunScaffold_Succeeds_WhenManifestFileExists(t *testing.T) {
	svc := scaffold.NewScaffoldService(&preventOverwriteStub{}, &manifestServiceErrorStub{}, &fileExistsStub{})

	testOut := testOut{sink: []string{}}
	result := svc.CreateScaffold(newScaffoldOptionsBuilder().withOut(&testOut).build())

	require.NoError(t, result)
	require.Len(t, testOut.sink, 1)
	assert.Contains(t, testOut.sink[0], "The Manifest file already exists, reusing:")
}

// ***************
// Test Stubs
// ***************

type preventOverwriteErrorStub struct{}

var errOverwriteError = errors.New("overwrite error")

func (*preventOverwriteErrorStub) PreventOverwrite(_, _ string, _ bool) error {
	return errOverwriteError
}

type preventOverwriteStub struct{}

func (*preventOverwriteStub) PreventOverwrite(_, _ string, _ bool) error {
	return nil
}

type manifestServiceErrorStub struct{}

func (*manifestServiceErrorStub) GetDefaultManifestContent() string {
	return "test"
}

func (*manifestServiceErrorStub) WriteManifestFile(_, _ string) error {
	return manifest.ErrWritingManifestFile
}

type manifestServiceStub struct{}

func (*manifestServiceStub) GetDefaultManifestContent() string {
	return "test"
}

func (*manifestServiceStub) WriteManifestFile(_, _ string) error {
	return nil
}

type fileExistsStub struct{}

func (*fileExistsStub) FileExists(_ string) (bool, error) {
	return true, nil
}

type fileDoesNotExistStub struct{}

func (*fileDoesNotExistStub) FileExists(_ string) (bool, error) {
	return false, nil
}

type testOut struct {
	sink []string
}

func (o *testOut) Write(msg string) {
	o.sink = append(o.sink, msg)
}

// ***************
// Test Options Builder
// ***************

type scaffoldOptionsBuilder struct {
	options scaffold.Options
}

func newScaffoldOptionsBuilder() *scaffoldOptionsBuilder {
	builder := &scaffoldOptionsBuilder{
		options: scaffold.Options{},
	}

	return builder.
		withOut(iotools.NewDefaultOut(io.Discard)).
		withDirectory("./").
		withModuleConfigFileName("scaffold-module-config.yaml").
		withManifestFileName("manifest.yaml").
		withModuleConfigFileOverwrite(false)
}

func (b *scaffoldOptionsBuilder) build() scaffold.Options {
	return b.options
}

func (b *scaffoldOptionsBuilder) withOut(out iotools.Out) *scaffoldOptionsBuilder {
	b.options.Out = out
	return b
}

func (b *scaffoldOptionsBuilder) withDirectory(directory string) *scaffoldOptionsBuilder {
	b.options.Directory = directory
	return b
}

func (b *scaffoldOptionsBuilder) withModuleConfigFileName(moduleConfigFileName string) *scaffoldOptionsBuilder {
	b.options.ModuleConfigFileName = moduleConfigFileName
	return b
}

func (b *scaffoldOptionsBuilder) withModuleConfigFileOverwrite(moduleConfigFileOverwrite bool) *scaffoldOptionsBuilder {
	b.options.ModuleConfigFileOverwrite = moduleConfigFileOverwrite
	return b
}

func (b *scaffoldOptionsBuilder) withManifestFileName(manifestFileName string) *scaffoldOptionsBuilder {
	b.options.ManifestFileName = manifestFileName
	return b
}
