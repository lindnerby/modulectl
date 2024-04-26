package scaffold

import (
	"fmt"

	"github.com/kyma-project/modulectl/tools/io"
)

type Options struct {
	Out                       io.Out
	Directory                 string
	ModuleConfigFileName      string
	ModuleConfigFileOverwrite bool
}

func (opts Options) validate() error {
	if opts.Out == nil {
		return fmt.Errorf("%w: opts.Out must not be nil", ErrInvalidOption)
	}

	if opts.Directory == "" {
		return fmt.Errorf("%w: opts.Directory must not be empty", ErrInvalidOption)
	}

	if opts.ModuleConfigFileName == "" {
		return fmt.Errorf("%w: opts.ModuleConfigFileName must not be empty", ErrInvalidOption)
	}

	return nil
}
