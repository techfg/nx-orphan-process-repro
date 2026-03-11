package main

import (
	"os"

	"github.com/techfg/nx-orphan-process-repro-cli/pkg/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
