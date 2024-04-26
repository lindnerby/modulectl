package scaffold

import (
	"fmt"

	"github.com/spf13/cobra"
)

func RunE(cmd *cobra.Command, args []string) error {
	// fmt.Fprintln(cmd.OutOrStdout(), "Validating")

	// if err := ValidateFlags(); err != nil {
	// 	return fmt.Errorf("%w\n", err)
	// }

	// if moduleConfigExists, err := ModuleConfigFileExists(); err != nil {
	// 	return fmt.Errorf("%w\n", err)
	// } else if moduleConfigExists && !Overwrite {
	// 	return fmt.Errorf("%w\n", errModuleConfigExists)
	// }

	// if manifestExists, err := ManifestFileExists(); err != nil {
	// 	return fmt.Errorf("%w\n", err)
	// } else if manifestExists {
	// 	fmt.Fprintf(cmd.OutOrStdout(), "The Manifest file already exists, reusing: %s\n", ManifestFilePath())
	// } else {
	// 	fmt.Fprintln(cmd.OutOrStdout(), "Generating Manifest file")
	// 	if err := GenerateManifest(); err != nil {
	// 		return fmt.Errorf("%w: %s\n", errManifestCreation, ManifestFilePath())
	// 	}

	// 	fmt.Fprintf(cmd.OutOrStdout(), "Generated a blank Manifest file: %s\n", ManifestFilePath())
	// }

	if defaultCRFileConfigured() {
		defaultCRFileExists, err := DefaultCRFileExists()
		if err != nil {
			return err
		}
		if !defaultCRFileExists {
			fmt.Fprintln(cmd.OutOrStdout(), "Generating default CR file")
			if err := GenerateDefaultCRFile(); err != nil {
				return fmt.Errorf("%w: %s\n", errDefaultCRCreationFailed, err.Error())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Generated a blank default CR file: %s\n", DefaultCRFilePath())
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "The default CR file already exists, reusing: %s\n", DefaultCRFilePath())
		}
	}

	if securityConfigFileConfigured() {
		secConfExists, err := SecurityConfigFileExists()
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Configuring Security-Scanners config file...")

		if !secConfExists {
			fmt.Fprintln(cmd.OutOrStdout(), "Generating security-scanners config file")
			if err := GenerateSecurityConfigFile(); err != nil {
				return fmt.Errorf("%w: %s\n", errSecurityConfigCreationFailed, err.Error())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Generated security-scanners config file: %s\n", SecurityConfigFilePath())
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "The security-scanners config file already exists, reusing: %s\n", SecurityConfigFilePath())
		}
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Generating module config file...")

	if err := GenerateModuleConfigFile(); err != nil {
		return fmt.Errorf("%w: %s\n", errModuleConfigCreationFailed, err.Error())
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Generated module config file: %s\n", ModuleConfigFilePath())

	return nil
}
