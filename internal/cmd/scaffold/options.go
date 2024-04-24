package scaffold

import (
	"fmt"

	"github.com/kyma-project/modulectl/tools/io"
)

type Options struct {
	Out io.Out
}

func (opts Options) validate() error {
	if opts.Out == nil {
		return fmt.Errorf("%w: opts.Out must not be nil", ErrInvalidOption)
	}

	return nil
}
