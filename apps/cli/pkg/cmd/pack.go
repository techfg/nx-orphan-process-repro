package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:          "pack",
		Short:        "cli pack",
		RunE:         packer,
		SilenceUsage: true,
	})
}

func packer(cmd *cobra.Command, args []string) error {
	if err := os.MkdirAll("dist", 0o755); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join("dist", "lib.js"), []byte("compiled"), 0o644)
}
