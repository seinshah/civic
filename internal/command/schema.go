package command

import (
	"github.com/seinshah/civic/internal/pkg/types"
	"github.com/seinshah/civic/internal/schema"
	"github.com/spf13/cobra"
)

func (c *Command) getConfigCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Helper commands to work with CV schema file.",
		Long:  `Use schema's sub-commands to help you manage the CV schema file.`,
	}

	cmd.AddCommand(getSchemaInitCommand())

	return cmd
}

func getSchemaInitCommand() *cobra.Command {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a new template CV schema file.",
		Long:  `A kickstart for you to create a sample schema file and start building your CV.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			handler := schema.NewHandler()

			return handler.Init(cmd.Context(), outputPath)
		},
	}

	cmd.Flags().StringVarP(
		&outputPath,
		"output", "o", types.CurrentWDPath(types.DefaultSchemaFileName),
		`path to the output configuration file.`,
	)

	return cmd
}
