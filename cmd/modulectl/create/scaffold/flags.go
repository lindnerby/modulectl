package scaffold

import (
	"github.com/kyma-project/modulectl/internal/scaffold"
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

	ManifestFileFlagName    = "gen-manifest"
	ManifestFileFlagDefault = "manifest.yaml"
	manifestFileFlagUsage   = "Specifies the manifest in the generated module config. A blank manifest file is generated if it doesn't exist"

	DefaultCRFlagName    = "gen-default-cr"
	DefaultCRFlagDefault = "default-cr.yaml"
	defaultCRFlagUsage   = "Specifies the default CR in the generated module config. A blank default CR file is generated if it doesn't exist"
)

func parseFlags(flags *pflag.FlagSet, opts *scaffold.Options) error {
	flags.StringVarP(&opts.Directory, DirectoryFlagName, directoryFlagShort, DirectoryFlagDefault, directoryFlagUsage)
	flags.StringVar(&opts.ModuleConfigFileName, ModuleConfigFileFlagName, ModuleConfigFileFlagDefault, moduleConfigFileFlagUsage)
	flags.BoolVarP(&opts.ModuleConfigFileOverwrite, ModuleConfigFileOverwriteFlagName, moduleConfigFileOverwriteFlagShort, ModuleConfigFileOverwriteFlagDefault, moduleConfigFileOverwriteFlagUsage)
	flags.StringVar(&opts.ManifestFileName, ManifestFileFlagName, ManifestFileFlagDefault, manifestFileFlagUsage)
	flags.StringVar(&opts.DefaultCRFileName, DefaultCRFlagName, DefaultCRFlagDefault, defaultCRFlagUsage)

	return nil
}
