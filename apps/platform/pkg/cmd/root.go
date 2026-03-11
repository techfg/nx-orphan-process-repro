package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "Platform",
	Short:         "The Platform",
	SilenceErrors: true,
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}
