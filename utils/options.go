package utils

import (
	"github.com/kyma-project/modulectl/utils/step"
)

// Options defines available options for the command
type Options struct {
	CI      bool
	Verbose bool
	step.Factory
	KubeconfigPath string
	Finalizers     *Finalizers
}

// NewOptions creates options with default values
func NewOptions() *Options {
	return &Options{
		Finalizers: NewFinalizer(),
	}
}
