package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/techfg/nx-orphan-process-repro-platform/pkg/cmd"
)

func main() {
	ctx := context.Background()
	if err := cmd.Execute(ctx); err != nil {
		slog.ErrorContext(ctx, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
