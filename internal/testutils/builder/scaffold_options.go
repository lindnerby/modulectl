package builder

import (
	"io"

	"github.com/kyma-project/modulectl/internal/cmd/scaffold"
	iotools "github.com/kyma-project/modulectl/tools/io"
)

type ScaffoldOptionsBuilder struct {
	options scaffold.Options
}

func NewScaffoldOptionsBuilder() *ScaffoldOptionsBuilder {
	builder := &ScaffoldOptionsBuilder{
		options: scaffold.Options{},
	}

	return builder.
		WithOut(iotools.NewDefaultOut(io.Discard)).
		WithDirectory("./").
		WithModuleConfigFileName("scaffold-module-config.yaml").
		WithModuleConfigFileOverwrite(false)
}

func (b *ScaffoldOptionsBuilder) Build() scaffold.Options {
	return b.options
}

func (b *ScaffoldOptionsBuilder) WithOut(out iotools.Out) *ScaffoldOptionsBuilder {
	b.options.Out = out
	return b
}

func (b *ScaffoldOptionsBuilder) WithDirectory(directory string) *ScaffoldOptionsBuilder {
	b.options.Directory = directory
	return b
}

func (b *ScaffoldOptionsBuilder) WithModuleConfigFileName(moduleConfigFileName string) *ScaffoldOptionsBuilder {
	b.options.ModuleConfigFileName = moduleConfigFileName
	return b
}

func (b *ScaffoldOptionsBuilder) WithModuleConfigFileOverwrite(moduleConfigFileOverwrite bool) *ScaffoldOptionsBuilder {
	b.options.ModuleConfigFileOverwrite = moduleConfigFileOverwrite
	return b
}
