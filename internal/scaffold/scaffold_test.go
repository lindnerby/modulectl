package scaffold_test

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kyma-project/modulectl/internal/scaffold"
	"github.com/kyma-project/modulectl/internal/scaffold/defaultcr"
	"github.com/kyma-project/modulectl/internal/scaffold/manifest"
	iotools "github.com/kyma-project/modulectl/tools/io"
)

func Test_RunScaffold_ReturnsError_WhenOutIsNil(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&preventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileDoesNotExistStub{})
	opts := newScaffoldOptionsBuilder().withOut(nil).build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.Out")
}

func Test_RunScaffold_ReturnsError_WhenDirectoryIsEmpty(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&preventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileDoesNotExistStub{})
	opts := newScaffoldOptionsBuilder().withDirectory("").build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.Directory")
}

func Test_RunScaffold_ReturnsError_WhenModuleConfigFileIsEmpty(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&preventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileDoesNotExistStub{})
	opts := newScaffoldOptionsBuilder().withModuleConfigFileName("").build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ModuleConfigFileName")
}

func Test_RunScaffold_ReturnsError_WhenManifestFileIsEmpty(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&preventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileDoesNotExistStub{})
	opts := newScaffoldOptionsBuilder().withManifestFileName("").build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ManifestFileName")
}

func Test_RunScaffold_ReturnsError_WhenModuleConfigServicePreventOverwriteReturnsError(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&preventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileDoesNotExistStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.ErrorIs(t, result, errOverwriteError)
}

func Test_RunScaffold_ReturnsError_WhenGeneratingManifestFileFails(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&preventOverwriteStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileDoesNotExistStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.ErrorIs(t, result, manifest.ErrGeneratingManifestFile)
}

func Test_RunScaffold_Succeeds_WhenGeneratingManifestFile(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&preventOverwriteStub{},
		&manifestServiceStub{},
		&defaultCRServiceStub{},
		&fileDoesNotExistStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.NoError(t, result)
}

func Test_RunScaffold_Succeeds_WhenDefaultCRFileIsNotConfigured(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&preventOverwriteStub{},
		&manifestServiceStub{},
		&defaultCRServiceErrorStub{},
		&fileDoesNotExistStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().withDefaultCRFileName("").build())

	require.NoError(t, result)
}

func Test_RunScaffold_ReturnsError_WhenGeneratingDefaultCRFileFails(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&preventOverwriteStub{},
		&manifestServiceStub{},
		&defaultCRServiceErrorStub{},
		&fileDoesNotExistStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.ErrorIs(t, result, defaultcr.ErrGeneratingDefaultCRFile)
}

func Test_RunScaffold_Succeeds_WhenGeneratingDefaultCRFile(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&preventOverwriteStub{},
		&manifestServiceStub{},
		&defaultCRServiceStub{},
		&fileDoesNotExistStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.NoError(t, result)
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

type manifestServiceStub struct{}

func (*manifestServiceStub) GenerateManifestFile(_ iotools.Out, _ string) error {
	return nil
}

type manifestServiceErrorStub struct{}

func (*manifestServiceErrorStub) GenerateManifestFile(_ iotools.Out, _ string) error {
	return manifest.ErrGeneratingManifestFile
}

type defaultCRServiceStub struct{}

func (*defaultCRServiceStub) GenerateDefaultCRFile(_ iotools.Out, _ string) error {
	return nil
}

type defaultCRServiceErrorStub struct{}

func (*defaultCRServiceErrorStub) GenerateDefaultCRFile(out iotools.Out, _ string) error {
	return defaultcr.ErrGeneratingDefaultCRFile
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
		withDefaultCRFileName("default-cr.yaml").
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

func (b *scaffoldOptionsBuilder) withDefaultCRFileName(defaultCRFileName string) *scaffoldOptionsBuilder {
	b.options.DefaultCRFileName = defaultCRFileName
	return b
}
