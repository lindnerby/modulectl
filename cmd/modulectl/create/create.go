package create

import (
	"github.com/kyma-project/modulectl/cmd/modulectl/create/scaffold"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates artifacts related to module development",
		Long:  `Use this command to create artifacts that are needed for module development.`,
	}

	cmd.AddCommand(scaffold.NewCmd())

	return cmd
}
