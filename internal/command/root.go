package command

import (
	"log/slog"
	"os"

	"github.com/seinshah/cvci/internal/cvhandler"
	"github.com/spf13/cobra"
)

// CommandExecutor is an interface that should be implemented by the command manager.
// The user of this interface only deals with the Execute method.
type CommandExecutor interface {
	Execute() error
}

// Command is the default command manager for this project.
// It is responsible to register all the necessary command and sub-commands.
type Command struct {
	root    *cobra.Command
	verbose bool
	debug   bool
}

var _ CommandExecutor = (*Command)(nil)

// NewCommand creates a new instance of the command manager.
// It is responsible to register all the necessary command and sub-commands.
// Register commands are:
// - generate.
func NewCommand(version string) *Command {
	var cmd Command

	cmd.root = &cobra.Command{
		Use:     "cvci",
		Version: version,
		Short:   "CVCI is a resume or cv as a code",
		Long: `A simple approach to help you maintain your resume or cv as a code
			and keep it in a git repository for easy and continuous updates.`,
	}

	cmd.root.PersistentFlags().BoolVar(
		&cmd.verbose, "verbose", false,
		"show more information during the process",
	)
	cmd.root.PersistentFlags().BoolVar(
		&cmd.debug, "debug", false,
		"enable the debug logs",
	)

	cmd.root.AddCommand(cmd.getGeneratorCommand())

	return &cmd
}

func (c *Command) Execute() error {
	return c.root.Execute()
}

func (c *Command) getGeneratorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate the resume or cv",
		Long:  `Generate the resume or cv based on the configuration file and the provided version.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c.updateLogLevel()

			cv := cvhandler.NewGenerator("")

			return cv.Generate()
		},
	}

	return cmd
}

func (c *Command) updateLogLevel() {
	// TODO: change it in Go version 1.22 to only update the default level
	if c.debug {
		slog.SetDefault(
			slog.New(
				slog.NewTextHandler(
					os.Stderr,
					&slog.HandlerOptions{
						Level: slog.LevelDebug,
					},
				),
			),
		)
	} else if c.verbose {
		slog.SetDefault(
			slog.New(
				slog.NewTextHandler(
					os.Stderr,
					&slog.HandlerOptions{
						Level: slog.LevelInfo,
					},
				),
			),
		)
	}
}
