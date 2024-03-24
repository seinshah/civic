package command

import (
	"context"
	"log/slog"

	"github.com/seinshah/cvci/internal/pkg/logger"
	"github.com/spf13/cobra"
)

// Command is the default command manager for this project.
// It is responsible to register all the necessary command and sub-commands.
type Command struct {
	root    *cobra.Command
	verbose bool
	debug   bool
	noColor bool
}

// NewCommand creates a new instance of the command manager.
// It is responsible to register all the necessary command and sub-commands.
func NewCommand(version string) *Command {
	var cmd Command

	cmd.root = &cobra.Command{
		Use:     "cvci",
		Version: version,
		Short:   "CVCI is a resume or cv as a code",
		Long:    `A simple approach to help you maintain your resume or cv as a code.`,
		PreRun: func(_ *cobra.Command, _ []string) {
			cmd.updateLogLevel()
		},
	}

	cmd.root.PersistentFlags().BoolVar(
		&cmd.verbose, "verbose", false,
		"show more information during the process",
	)
	cmd.root.PersistentFlags().BoolVar(
		&cmd.debug, "debug", false,
		"enable the debug logs",
	)
	cmd.root.PersistentFlags().BoolVar(
		&cmd.noColor, "no-color", false,
		"disable the color in the logs output",
	)

	cmd.root.AddCommand(
		cmd.getGenerateCommands(),
		cmd.getConfigCommands(),
	)

	return &cmd
}

func (c *Command) Execute(ctx context.Context) error {
	return c.root.ExecuteContext(ctx)
}

func (c *Command) updateLogLevel() {
	opts := make([]logger.Option, 0)

	if c.debug {
		opts = append(opts, logger.WithLevel(slog.LevelDebug))
	} else if c.verbose {
		opts = append(opts, logger.WithLevel(slog.LevelInfo))
	}

	if c.noColor {
		opts = append(opts, logger.WithNoColor())
	}

	logger.SetUp(opts...)
}
