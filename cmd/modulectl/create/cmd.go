package create

import (
	_ "embed"

	"github.com/kyma-project/modulectl/cmd/modulectl/create/scaffold"
	"github.com/spf13/cobra"
)

//go:embed use.txt
var use string

//go:embed short.txt
var short string

//go:embed long.txt
var long string

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(scaffold.NewCmd())

	return cmd
}
