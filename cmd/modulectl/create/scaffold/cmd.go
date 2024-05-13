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

//go:embed use.txt
var use string

//go:embed short.txt
var short string

//go:embed long.txt
var long string

//go:embed example.txt
var example string

func NewCmd() *cobra.Command {
	fileSystemUtil := &filesystem.FileSystemUtil{}
	yamlConverter := &yaml.ObjectToYAMLConverter{}

	moduleConfigContentProvider, err := contentprovider.NewModuleConfigContentProvider(yamlConverter)
	if err != nil {
		panic(err)
	}

	moduleConfigFileGenerator, err := filegenerator.NewFileGeneratorService("module-config", fileSystemUtil, moduleConfigContentProvider)
	if err != nil {
		panic(err)
	}

	moduleConfigService, err := moduleconfig.NewModuleConfigService(fileSystemUtil, moduleConfigFileGenerator)
	if err != nil {
		panic(err)
	}

	manifestFileGenerator, err := filegenerator.NewFileGeneratorService("manifest", fileSystemUtil, contentprovider.NewManifestContentProvider())
	if err != nil {
		panic(err)
	}

	manifestReuseFileGenerator, err := filegenerator.NewReuseFileGeneratorService("manifest", fileSystemUtil, manifestFileGenerator)
	if err != nil {
		panic(err)
	}

	defaultCRFileGenerator, err := filegenerator.NewFileGeneratorService("defaultcr", fileSystemUtil, contentprovider.NewDefaultCRContentProvider())
	if err != nil {
		panic(err)
	}

	defaultCRReuseFileGenerator, err := filegenerator.NewReuseFileGeneratorService("defaultcr", fileSystemUtil, defaultCRFileGenerator)
	if err != nil {
		panic(err)
	}

	securitConfigContentProvider, err := contentprovider.NewSecurityConfigContentProvider(yamlConverter)
	if err != nil {
		panic(err)
	}

	securityConfigFileGenerator, err := filegenerator.NewFileGeneratorService("security-config", fileSystemUtil, securitConfigContentProvider)
	if err != nil {
		panic(err)
	}

	securityConfigReuseFileGenerator, err := filegenerator.NewReuseFileGeneratorService("security-config", fileSystemUtil, securityConfigFileGenerator)
	if err != nil {
		panic(err)
	}

	scaffoldService, err := scaffold.NewScaffoldService(
		moduleConfigService,
		manifestReuseFileGenerator,
		defaultCRReuseFileGenerator,
		securityConfigReuseFileGenerator,
	)
	if err != nil {
		panic(err)
	}

	opts := scaffold.Options{}

	cmd := &cobra.Command{
		Use:     use,
		Short:   short,
		Long:    long,
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
