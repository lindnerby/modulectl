package scaffold_test

import (
	"errors"
	"io"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kyma-project/modulectl/internal/scaffold"
	"github.com/kyma-project/modulectl/internal/scaffold/common/types"
	"github.com/kyma-project/modulectl/internal/scaffold/defaultcr"
	"github.com/kyma-project/modulectl/internal/scaffold/manifest"
	iotools "github.com/kyma-project/modulectl/tools/io"
)

func Test_RunScaffold_ReturnsError_WhenOutIsNil(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigPreventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})
	opts := newScaffoldOptionsBuilder().withOut(nil).build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.Out")
}

func Test_RunScaffold_ReturnsError_WhenDirectoryIsEmpty(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigPreventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})
	opts := newScaffoldOptionsBuilder().withDirectory("").build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.Directory")
}

func Test_RunScaffold_ReturnsError_WhenModuleConfigFileIsEmpty(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigPreventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})
	opts := newScaffoldOptionsBuilder().withModuleConfigFileName("").build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ModuleConfigFileName")
}

func Test_RunScaffold_ReturnsError_WhenManifestFileIsEmpty(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigPreventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})
	opts := newScaffoldOptionsBuilder().withManifestFileName("").build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ManifestFileName")
}

func Test_RunScaffold_ReturnsError_WhenModuleNameIsEmpty(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigPreventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})
	opts := newScaffoldOptionsBuilder().withModuleName("").build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ModuleName")
}

func Test_RunScaffold_ReturnsError_WhenModuleNameIsExceedingLength(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigPreventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})
	opts := newScaffoldOptionsBuilder().withModuleName(getRandomName(256)).build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ModuleName")
	assert.Contains(t, result.Error(), "length")
}

func Test_RunScaffold_ReturnsError_WhenModuleNameIsNotMatchingPattern(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigPreventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})
	opts := newScaffoldOptionsBuilder().withModuleName(getRandomName(10)).build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ModuleName")
	assert.Contains(t, result.Error(), "pattern")
}

func Test_RunScaffold_ReturnsError_WhenModuleVersionIsEmpty(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigPreventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})
	opts := newScaffoldOptionsBuilder().withModuleVersion("").build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ModuleVersion")
}

func Test_RunScaffold_ReturnsError_WhenModuleVersionIsInvalid(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigPreventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})
	opts := newScaffoldOptionsBuilder().withModuleVersion(getRandomName(10)).build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ModuleVersion")
	assert.Contains(t, result.Error(), "failed to parse")
}

func Test_RunScaffold_ReturnsError_WhenModuleChannelIsEmpty(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigPreventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})
	opts := newScaffoldOptionsBuilder().withModuleChannel("").build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ModuleChannel")
}

func Test_RunScaffold_ReturnsError_WhenModuleChannelIsExceedingLength(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigPreventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})
	opts := newScaffoldOptionsBuilder().withModuleChannel(getRandomName(33)).build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ModuleChannel")
	assert.Contains(t, result.Error(), "length")
}

func Test_RunScaffold_ReturnsError_WhenModuleChannelFallsBelowLength(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigPreventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})
	opts := newScaffoldOptionsBuilder().withModuleChannel(getRandomName(2)).build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ModuleChannel")
	assert.Contains(t, result.Error(), "length")
}

func Test_RunScaffold_ReturnsError_WhenModuleChannelNotMatchingCharset(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigPreventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})
	opts := newScaffoldOptionsBuilder().withModuleChannel("with not allowed chars 123%&").build()

	result := svc.CreateScaffold(opts)

	require.ErrorIs(t, result, scaffold.ErrInvalidOption)
	assert.Contains(t, result.Error(), "opts.ModuleChannel")
	assert.Contains(t, result.Error(), "pattern")
}

func Test_RunScaffold_ReturnsError_WhenModuleConfigServicePreventOverwriteReturnsError(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigPreventOverwriteErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.ErrorIs(t, result, errOverwriteError)
}

func Test_RunScaffold_ReturnsError_WhenGeneratingManifestFileFails(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigGenerateFileErrorStub{},
		&manifestServiceErrorStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.ErrorIs(t, result, manifest.ErrGeneratingManifestFile)
}

func Test_RunScaffold_Succeeds_WhenGeneratingManifestFile(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigStub{},
		&manifestServiceStub{},
		&defaultCRServiceStub{},
		&fileGeneratorStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.NoError(t, result)
}

func Test_RunScaffold_Succeeds_WhenDefaultCRFileIsNotConfigured(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigStub{},
		&manifestServiceStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().withDefaultCRFileName("").build())

	require.NoError(t, result)
}

func Test_RunScaffold_ReturnsError_WhenGeneratingDefaultCRFileFails(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigGenerateFileErrorStub{},
		&manifestServiceStub{},
		&defaultCRServiceErrorStub{},
		&fileGeneratorErrorStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.ErrorIs(t, result, defaultcr.ErrGeneratingDefaultCRFile)
}

func Test_RunScaffold_Succeeds_WhenGeneratingDefaultCRFile(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigStub{},
		&manifestServiceStub{},
		&defaultCRServiceStub{},
		&fileGeneratorStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.NoError(t, result)
}

func Test_RunScaffold_Succeeds_WhenSecurityConfigFileIsNotConfigured(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigStub{},
		&manifestServiceStub{},
		&defaultCRServiceStub{},
		&fileGeneratorErrorStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().withSecurityConfigFileName("").build())

	require.NoError(t, result)
}

func Test_RunScaffold_ReturnsError_WhenGeneratingSecurityConfigFileFails(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigGenerateFileErrorStub{},
		&manifestServiceStub{},
		&defaultCRServiceStub{},
		&fileGeneratorErrorStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.ErrorIs(t, result, scaffold.ErrGenertingFile)
	require.ErrorIs(t, result, errSomeFileGeneratorError)
	assert.Contains(t, result.Error(), "security-config.yaml")
}

func Test_RunScaffold_Succeeds_WhenGeneratingSecurityConfigFile(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigStub{},
		&manifestServiceStub{},
		&defaultCRServiceStub{},
		&fileGeneratorStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.NoError(t, result)
}

func Test_RunScaffold_ReturnsError_WhenGeneratingModuleConfigReturnsError(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigGenerateFileErrorStub{},
		&manifestServiceStub{},
		&defaultCRServiceStub{},
		&fileGeneratorStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.ErrorIs(t, result, scaffold.ErrGenertingFile)
	require.ErrorIs(t, result, errSomeFileGeneratorError)
	assert.Contains(t, result.Error(), "module-config.yaml")
}

func Test_RunScaffold_Succeeds(t *testing.T) {
	svc := scaffold.NewScaffoldService(
		&moduleConfigStub{},
		&manifestServiceStub{},
		&defaultCRServiceStub{},
		&fileGeneratorStub{})

	result := svc.CreateScaffold(newScaffoldOptionsBuilder().build())

	require.NoError(t, result)
}

// ***************
// Test Stubs
// ***************
var errSomeFileGeneratorError = errors.New("some file generator error")
var errSomeUnexpectedError = errors.New("if you see this error, something went wrong in the test setup")

type moduleConfigPreventOverwriteErrorStub struct{}

var errOverwriteError = errors.New("overwrite error")

func (*moduleConfigPreventOverwriteErrorStub) PreventOverwrite(_, _ string, _ bool) error {
	return errOverwriteError
}

func (*moduleConfigPreventOverwriteErrorStub) GenerateFile(_ iotools.Out, _ string, _ types.KeyValueArgs) error {
	return errSomeUnexpectedError
}

type moduleConfigGenerateFileErrorStub struct{}

func (*moduleConfigGenerateFileErrorStub) PreventOverwrite(_, _ string, _ bool) error {
	return nil
}

func (*moduleConfigGenerateFileErrorStub) GenerateFile(_ iotools.Out, _ string, _ types.KeyValueArgs) error {
	return errSomeFileGeneratorError
}

type moduleConfigStub struct{}

func (*moduleConfigStub) PreventOverwrite(_, _ string, _ bool) error {
	return nil
}

func (*moduleConfigStub) GenerateFile(_ iotools.Out, _ string, _ types.KeyValueArgs) error {
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

type fileGeneratorErrorStub struct{}

func (*fileGeneratorErrorStub) GenerateFile(_ iotools.Out, _ string, _ types.KeyValueArgs) error {
	return errSomeFileGeneratorError
}

type fileGeneratorStub struct{}

func (*fileGeneratorStub) GenerateFile(_ iotools.Out, _ string, _ types.KeyValueArgs) error {
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
		withManifestFileName("manifest.yaml").
		withDefaultCRFileName("default-cr.yaml").
		withModuleConfigFileOverwrite(false).
		withSecurityConfigFileName("security-config.yaml").
		withModuleName("github.com/kyma-project/test").
		withModuleVersion("0.0.1").
		withModuleChannel("experimental")
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

func (b *scaffoldOptionsBuilder) withSecurityConfigFileName(securityConfigFileName string) *scaffoldOptionsBuilder {
	b.options.SecurityConfigFileName = securityConfigFileName
	return b
}

func (b *scaffoldOptionsBuilder) withModuleName(moduleName string) *scaffoldOptionsBuilder {
	b.options.ModuleName = moduleName
	return b
}

func (b *scaffoldOptionsBuilder) withModuleVersion(moduleVersion string) *scaffoldOptionsBuilder {
	b.options.ModuleVersion = moduleVersion
	return b
}

func (b *scaffoldOptionsBuilder) withModuleChannel(moduleChannel string) *scaffoldOptionsBuilder {
	b.options.ModuleChannel = moduleChannel
	return b
}

// ***************
// Test Helpers
// ***************

const charset = "abcdefghijklmnopqrstuvwxyz"

func getRandomName(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
