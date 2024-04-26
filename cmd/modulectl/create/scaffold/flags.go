package scaffold

import (
	"github.com/kyma-project/modulectl/internal/cmd/scaffold"
	"github.com/spf13/pflag"
)

const (
	DirectoryFlagName    = "directory"
	directoryFlagShort   = "d"
	DirectoryFlagDefault = "./"
	directoryFlagUsage   = "Specifies the target directory where the scaffolding shall be generated"

	ModuleConfigFileFlagName    = "module-config"
	ModuleConfigFileFlagDefault = "scaffold-module-config.yaml"
	moduleConfigFileFlagUsage   = "Specifies the name of the generated module configuration file"

	ModuleConfigFileOverwriteFlagName    = "overwrite"
	moduleConfigFileOverwriteFlagShort   = "o"
	ModuleConfigFileOverwriteFlagDefault = false
	moduleConfigFileOverwriteFlagUsage   = "Specifies if the command overwrites an existing module configuration file"
)

func parseFlags(flags *pflag.FlagSet, opts *scaffold.Options) error {
	flags.StringVarP(&opts.Directory, DirectoryFlagName, directoryFlagShort, DirectoryFlagDefault, directoryFlagUsage)
	flags.StringVar(&opts.ModuleConfigFileName, ModuleConfigFileFlagName, ModuleConfigFileFlagDefault, moduleConfigFileFlagUsage)
	flags.BoolVarP(&opts.ModuleConfigFileOverwrite, ModuleConfigFileOverwriteFlagName, moduleConfigFileOverwriteFlagShort, ModuleConfigFileOverwriteFlagDefault, moduleConfigFileOverwriteFlagUsage)

	return nil
}
