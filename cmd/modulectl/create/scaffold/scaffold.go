package scaffold

import (
	_ "embed"

	"github.com/spf13/cobra"

	"github.com/kyma-project/modulectl/internal/scaffold"
	"github.com/kyma-project/modulectl/internal/scaffold/contentprovider"
	"github.com/kyma-project/modulectl/internal/scaffold/filegenerator"
	"github.com/kyma-project/modulectl/internal/scaffold/moduleconfig"
	"github.com/kyma-project/modulectl/tools/filesystem"
	"github.com/kyma-project/modulectl/tools/io"
	"github.com/kyma-project/modulectl/tools/yaml"
)

//go:embed cmd_description.txt
var description string

//go:embed cmd_example.txt
var example string

func NewCmd() *cobra.Command {
	fileSystemUtil := &filesystem.FileSystemUtil{}
	yamlConverter := &yaml.ObjectToYAMLConverter{}
	scaffoldService := scaffold.NewScaffoldService(
		moduleconfig.NewModuleConfigService(
			fileSystemUtil,
			filegenerator.NewFileGeneratorService(
				"module-config",
				fileSystemUtil,
				contentprovider.NewModuleConfigContentProvider(yamlConverter),
			),
		),
		filegenerator.NewReuseFileGeneratorService(
			"manifest",
			fileSystemUtil,
			filegenerator.NewFileGeneratorService(
				"manifest",
				fileSystemUtil,
				contentprovider.NewManifestContentProvider(),
			),
		),
		filegenerator.NewReuseFileGeneratorService(
			"defaultcr",
			fileSystemUtil,
			filegenerator.NewFileGeneratorService(
				"defaultcr",
				fileSystemUtil,
				contentprovider.NewDefaultCRContentProvider(),
			),
		),
		filegenerator.NewReuseFileGeneratorService(
			"security-config",
			fileSystemUtil,
			filegenerator.NewFileGeneratorService(
				"security-config",
				fileSystemUtil,
				contentprovider.NewSecurityConfigContentProvider(yamlConverter))),
	)

	opts := scaffold.Options{}

	cmd := &cobra.Command{
		Use:     "scaffold [--module-name MODULE_NAME --module-version MODULE_VERSION --module-channel CHANNEL] [--directory MODULE_DIRECTORY] [flags]",
		Short:   "Generates necessary files required for module creation",
		Long:    description,
		Example: example,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return scaffoldService.CreateScaffold(opts)
		},
	}

	opts.Out = io.NewDefaultOut(cmd.OutOrStdout())
	parseFlags(cmd.Flags(), &opts)

	return cmd
}
