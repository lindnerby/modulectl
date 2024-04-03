package modulectl

import (
	"fmt"
	"github.com/kyma-project/modulectl/cmd/modulectl/create"
	"github.com/kyma-project/modulectl/utils"
	"github.com/spf13/cobra"
)

func NewCmd(o *utils.Options) *cobra.Command {
	rootCommand := &cobra.Command{
		Use:   "modulectl",
		Short: "This is the Kyma Module Controller CLI",
		Long:  "A CLI from the Kyma Module Controller. Wonderful to use.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("This is the Kyma ModuleCtl command executed")
		},
	}

	rootCommand.AddCommand(create.NewCmd(o))

	return rootCommand
}
