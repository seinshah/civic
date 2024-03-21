package command

import (
	"context"
	"log/slog"

	"github.com/seinshah/cvci/internal/cvhandler"
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
	cmd.root.PersistentFlags().BoolVar(
		&cmd.noColor, "no-color", false,
		"disable the color in the logs output",
	)

	cmd.root.AddCommand(cmd.getGeneratorCommand())

	return &cmd
}

func (c *Command) Execute(ctx context.Context) error {
	return c.root.ExecuteContext(ctx)
}

func (c *Command) getGeneratorCommand() *cobra.Command {
	var (
		configPath string
		outputPath string
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate the resume or cv",
		Long:  `Generate the resume or cv based on the configuration file and the provided version.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c.updateLogLevel()

			cv, err := cvhandler.NewGenerator(cmd.Version, configPath, outputPath)
			if err != nil {
				return err
			}

			return cv.Generate(cmd.Context())
		},
	}

	cmd.Flags().StringVarP(
		&configPath,
		"config", "c", "",
		"absolute path/link to the configuration file.",
	)

	cmd.Flags().StringVarP(
		&outputPath,
		"output", "o", "",
		"path to the output pdf file.",
	)

	return cmd
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
