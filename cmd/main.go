package main

import (
	"github.com/kyma-project/modulectl/cmd/modulectl"
	"github.com/kyma-project/modulectl/utils"
	"os"
)

func main() {
	command := modulectl.NewCmd(utils.NewOptions())

	if err := command.Execute(); err != nil {
		os.Exit(-1)
	}
}
