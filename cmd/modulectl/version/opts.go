package version

import (
	"github.com/kyma-project/modulectl/utils"
)

// Options defines available options for the command
type Options struct {
	*utils.Options
	ClientOnly bool
}

// NewOptions creates options with default values
func NewOptions(o *utils.Options) *Options {
	return &Options{Options: o}
}
