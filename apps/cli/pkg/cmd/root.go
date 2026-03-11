package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "CLI",
	Short: "The CLI",
}

func Execute() error {
	return rootCmd.Execute()
}
