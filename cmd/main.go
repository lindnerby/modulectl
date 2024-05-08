package main

import (
	"os"

	"github.com/kyma-project/modulectl/cmd/modulectl"
)

func main() {
	cmd := modulectl.NewCmd()

	if err := cmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
