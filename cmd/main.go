package main

import (
	"github.com/kyma-project/modulectl/cmd/modulectl"
	"os"
)

func main() {
	command := modulectl.NewCmd()

	if err := command.Execute(); err != nil {
		os.Exit(-1)
	}
}
