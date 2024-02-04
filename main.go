package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/seinshah/cvci/internal/cmd"
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

	command := cmd.NewCommand(Version)

	if err := command.Execute(); err != nil {
		slog.Error(fmt.Sprintf("failed to proceed: %s", err))
		os.Exit(1)
	}
}
