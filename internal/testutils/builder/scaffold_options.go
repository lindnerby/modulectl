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

	return builder.WithOut(iotools.NewDefaultOut(io.Discard))

}

func (b *ScaffoldOptionsBuilder) Build() scaffold.Options {
	return b.options
}

func (b *ScaffoldOptionsBuilder) WithOut(out iotools.Out) *ScaffoldOptionsBuilder {
	b.options.Out = out
	return b
}
