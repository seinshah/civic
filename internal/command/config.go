package command

import (
	"github.com/seinshah/cvci/internal/config"
	"github.com/seinshah/cvci/internal/pkg/types"
	"github.com/spf13/cobra"
)

func (c *Command) getConfigCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Helper commands to work with CVCI configuration file.",
		Long:  `Use config's sub-commands to help you manage the CVCI configuration file.`,
	}

	cmd.AddCommand(getConfigInitCommand())

	return cmd
}

func getConfigInitCommand() *cobra.Command {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a new template CVCI configuration file.",
		Long:  `Generate a new CVCI configuration file with sample values for all the sections.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			handler := config.NewHandler()

			return handler.Init(cmd.Context(), outputPath)
		},
	}

	cmd.Flags().StringVarP(
		&outputPath,
		"output", "o", "./"+types.DefaultConfigFileName,
		`path to the output configuration file. If not provided, it will be created in the current directory.`,
	)

	return cmd
}
