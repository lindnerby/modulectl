package scaffold

import (
	"github.com/spf13/cobra"

	"github.com/kyma-project/modulectl/internal/cmd/scaffold"
	"github.com/kyma-project/modulectl/internal/cmd/scaffold/moduleconfig"
	"github.com/kyma-project/modulectl/tools/filesystem"
	"github.com/kyma-project/modulectl/tools/io"
)

func NewCmd() *cobra.Command {
	scaffoldService := scaffold.NewScaffoldService(
		moduleconfig.NewModuleConfigService(
			&filesystem.FileSystemUtil{},
		),
	)

	opts := scaffold.Options{}

	cmd := &cobra.Command{
		Use:   "scaffold [--module-name MODULE_NAME --module-version MODULE_VERSION --module-channel CHANNEL] [--directory MODULE_DIRECTORY] [flags]",
		Short: "Generates necessary files required for module creation",
		Long: `Scaffold generates or configures the necessary files for creating a new module in Kyma. This includes setting up 
a basic directory structure and creating default files based on the provided flags.

The command is designed to streamline the module creation process in Kyma, making it easier and more 
efficient for developers to get started with new modules. It supports customization through various flags, 
allowing for a tailored scaffolding experience according to the specific needs of the module being created.

The command generates or uses the following files:
 - Module Config:
	Enabled: Always
	Adjustable with flag: --module-config=VALUE
	Generated when: The file doesn't exist or the --overwrite=true flag is provided
	Default file name: scaffold-module-config.yaml
 - Manifest:
	Enabled: Always
	Adjustable with flag: --gen-manifest=VALUE
	Generated when: The file doesn't exist. If the file exists, it's name is used in the generated module configuration file
	Default file name: manifest.yaml
 - Default CR(s):
	Enabled: When the flag --gen-default-cr is provided with or without value
	Adjustable with flag: --gen-default-cr[=VALUE], if provided without an explicit VALUE, the default value is used
	Generated when: The file doesn't exist. If the file exists, it's name is used in the generated module configuration file
	Default file name: default-cr.yaml
 - Security Scanners Config:
	Enabled: When the flag --gen-security-config is provided with or without value
	Adjustable with flag: --gen-security-config[=VALUE], if provided without an explicit VALUE, the default value is used
	Generated when: The file doesn't exist. If the file exists, it's name is used in the generated module configuration file
	Default file name: sec-scanners-config.yaml

**NOTE:**: To protect the user from accidental file overwrites, this command by default doesn't overwrite any files.
Only the module config file may be force-overwritten when the --overwrite=true flag is used.

You can specify the required fields of the module config using the following CLI flags:
--module-name=NAME
--module-version=VERSION
--module-channel=CHANNEL

**NOTE:**: If the required fields aren't provided, the defaults are applied and the module-config.yaml is not ready to be used. You must manually edit the file to make it usable.
Also, edit the sec-scanners-config.yaml to be able to use it.
`,
		Example: `Generate a minimal scaffold for a module - only a blank manifest file and module config file is generated using defaults
                modulectl create scaffold
Generate a scaffold providing required values explicitly
				modulectl create scaffold --module-name="kyma-project.io/module/testmodule" --module-version="0.1.1" --module-channel=fast
Generate a scaffold with a manifest file, default CR and security-scanners config for a module
				modulectl create scaffold --gen-default-cr --gen-security-config
Generate a scaffold with a manifest file, default CR and security-scanners config for a module, overriding default values
				modulectl create scaffold --gen-manifest="my-manifest.yaml" --gen-default-cr="my-cr.yaml" --gen-security-config="my-seccfg.yaml"

`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return scaffoldService.CreateScaffold(opts)
		},
	}

	opts.Out = io.NewDefaultOut(cmd.OutOrStdout())
	parseFlags(cmd.Flags(), &opts)

	// cmd.Flags().StringVar(&scaffold.ModuleName, "module-name", "kyma-project.io/module/mymodule",
	// 	"Specifies the module name in the generated config file")
	// cmd.Flags().StringVar(&scaffold.ModuleVersion, "module-version", "0.0.1",
	// 	"Specifies the module version in the generated module config file")
	// cmd.Flags().StringVar(&scaffold.ModuleChannel, "module-channel", "regular",
	// 	"Specifies the module channel in the generated module config file")

	// cmd.Flags().Lookup(scaffold.ModuleConfigFileFlagName).NoOptDefVal = scaffold.ModuleConfigFileFlagDefault

	// cmd.Flags().StringVar(&scaffold.ManifestFile, scaffold.ManifestFileFlagName, scaffold.ManifestFileFlagDefault,
	// 	"Specifies the manifest in the generated module config. A blank manifest file is generated if it doesn't exist")
	// cmd.Flags().Lookup(scaffold.ManifestFileFlagName).NoOptDefVal = scaffold.ManifestFileFlagDefault

	// cmd.Flags().StringVar(&scaffold.SecurityConfigFile, scaffold.SecurityConfigFlagName, "",
	// 	"Specifies the security file in the generated module config. A scaffold security config file is generated if it doesn't exist")
	// cmd.Flags().Lookup(scaffold.SecurityConfigFlagName).NoOptDefVal = scaffold.SecurityConfigFlagDefault

	// cmd.Flags().StringVar(&scaffold.DefaultCRFile, scaffold.DefaultCRFlagName, "",
	// 	"Specifies the defaultCR in the generated module config. A blank defaultCR file is generated if it doesn't exist")
	// cmd.Flags().Lookup(scaffold.DefaultCRFlagName).NoOptDefVal = scaffold.DefaultCRFlagDefault

	return cmd
}
