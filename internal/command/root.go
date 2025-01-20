package command

import (
	"context"
	"github.com/seinshah/civic/internal/pkg/types"
	"log/slog"

	"github.com/seinshah/civic/internal/pkg/logger"
	"github.com/spf13/cobra"
)

// Command is the default command manager for this project.
// It is responsible to register all the necessary command and sub-commands.
type Command struct {
	root    *cobra.Command
	verbose bool
	noColor bool
	version string
}

// NewCommand creates a new instance of the command manager.
// It is responsible to register all the necessary command and sub-commands.
func NewCommand(version string) *Command {
	cmd := Command{
		version: version,
	}

	cmd.root = &cobra.Command{
		Use:     types.DefaultAppName,
		Version: version,
		Short:   types.DefaultAppName + " - CV as a Code",
		Long: `A tool to help maintaining and extending CVs or resumes 
easily by separating the template from the content.`,
		PreRun: func(_ *cobra.Command, _ []string) {
			cmd.updateLogLevel()
		},
	}

	cmd.root.PersistentFlags().BoolVar(
		&cmd.verbose, "verbose", false,
		"show more information during the process",
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
	opts := []logger.Option{logger.WithLevel(slog.LevelDebug)}

	if c.noColor {
		opts = append(opts, logger.WithNoColor())
	}

	logger.SetUp(opts...)
}
