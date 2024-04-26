package scaffold_test

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kyma-project/modulectl/internal/scaffold"
	iotools "github.com/kyma-project/modulectl/tools/io"
)

func Test_RunScaffold_ReturnsError_WhenOutIsNil(t *testing.T) {
	svc := scaffold.NewScaffoldService(&preventOverwriteErrorStub{})
	opts := newScaffoldOptionsBuilder().withOut(nil).build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.Out")
}

func Test_RunScaffold_ReturnsError_WhenDirectoryIsEmpty(t *testing.T) {
	svc := scaffold.NewScaffoldService(&preventOverwriteErrorStub{})
	opts := newScaffoldOptionsBuilder().withDirectory("").build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.Directory")
}

func Test_RunScaffold_ReturnsError_WhenModuleConfigFileIsEmpty(t *testing.T) {
	svc := scaffold.NewScaffoldService(&preventOverwriteErrorStub{})
	opts := newScaffoldOptionsBuilder().withModuleConfigFileName("").build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ModuleConfigFileName")
}

func Test_RunScaffold_ReturnsError_WhenModuleConfigServicePreventOverwriteReturnsError(t *testing.T) {
	svc := scaffold.NewScaffoldService(&preventOverwriteErrorStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.ErrorIs(t, result, errOverwriteError)
}

func Test_RunScaffold_Succeeds(t *testing.T) {
	svc := scaffold.NewScaffoldService(&preventOverwriteStub{})

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
