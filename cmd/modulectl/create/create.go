package create

import (
	"github.com/kyma-project/modulectl/cmd/modulectl/create/scaffold"
	"github.com/kyma-project/modulectl/utils"
	"github.com/spf13/cobra"
)

// NewCmd creates a new Kyma CLI command
func NewCmd(o *utils.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates resources on the Kyma cluster.",
		Long: `Use this command to create resources on the Kyma cluster.
`,
	}

	cmd.AddCommand(scaffold.NewCmd(scaffold.NewOptions(o)))

	return cmd
}
