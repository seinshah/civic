package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/seinshah/cvci/internal/command"
)

// Version will be changed at build time using ldflags.
// The new value will be determined in the CI pipeline
// using the git tag.
var Version = "development"

func main() {
	// Setting up default logger with warning level
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(
				os.Stderr,
				&slog.HandlerOptions{
					Level: slog.LevelWarn,
				},
			),
		),
	)

	cmd := command.NewCommand(Version)

	if err := cmd.Execute(); err != nil {
		slog.Error(fmt.Sprintf("failed to proceed: %s", err))
		os.Exit(1)
	}
}
