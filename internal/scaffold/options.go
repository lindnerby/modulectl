package scaffold

import (
	"errors"
	"fmt"
	"os"

	"github.com/kyma-project/modulectl/tools/io"
)

type Options struct {
	Out                       io.Out
	Directory                 string
	ModuleConfigFileName      string
	ModuleConfigFileOverwrite bool
	ManifestFileName          string
	DefaultCRFileName         string
	SecurityConfigFileName    string
	ModuleName                string
	ModuleVersion             string
	ModuleChannel             string
}

func (opts Options) validate() error {
	if opts.Out == nil {
		return fmt.Errorf("%w: opts.Out must not be nil", ErrInvalidOption)
	}

	if err := opts.validateDirectory(); err != nil {
		return err
	}

	if opts.ModuleConfigFileName == "" {
		return fmt.Errorf("%w: opts.ModuleConfigFileName must not be empty", ErrInvalidOption)
	}

	if opts.ManifestFileName == "" {
		return fmt.Errorf("%w: opts.ManifestFileName must not be empty", ErrInvalidOption)
	}

	if opts.ModuleName == "" {
		return fmt.Errorf("%w: opts.ModuleName must not be empty", ErrInvalidOption)
	}

	if opts.ModuleVersion == "" {
		return fmt.Errorf("%w: opts.ModuleVersion must not be empty", ErrInvalidOption)
	}

	if opts.ModuleChannel == "" {
		return fmt.Errorf("%w: opts.ModuleChannel must not be empty", ErrInvalidOption)
	}

	return nil
}

// TODO check how this can be unit tested without making a service out of it
func (opts Options) validateDirectory() error {
	fileInfo, err := os.Stat(opts.Directory)

	if opts.Directory == "" {
		return fmt.Errorf("%w: opts.Directory must not be empty", ErrInvalidOption)
	}

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%w: directory %s does not exist", ErrInvalidOption, opts.Directory)
		}
		return err
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("%w: %s is not a directory", ErrInvalidOption, opts.Directory)
	}

	return nil
}

func (opts Options) defaultCRFileNameConfigured() bool {
	return opts.DefaultCRFileName != ""
}

func (opts Options) securityConfigFileNameConfigured() bool {
	return opts.SecurityConfigFileName != ""
}
