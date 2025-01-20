package command

import (
	"fmt"

	"github.com/seinshah/civic/internal/cv"
	"github.com/seinshah/civic/internal/pkg/types"
	"github.com/spf13/cobra"
)

func (c *Command) getGenerateCommands() *cobra.Command {
	var (
		schemaFilePath string
		outputPath     string
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate the resume or cv",
		Long:  `Generate the resume or cv based on the schema file and the provided version.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			handler, err := cv.NewHandler(c.version, schemaFilePath, outputPath)
			if err != nil {
				return err
			}

			return handler.Generate(cmd.Context())
		},
	}

	cmd.Flags().StringVarP(
		&schemaFilePath,
		"schema_file", "s", types.CurrentWDPath(types.DefaultSchemaFileName),
		`The local path or link to the CV schema file. This file includes the data that will be replicated in
the template. valid types: `+fmt.Sprintf("%v", types.SchemaTypeNames()),
	)

	cmd.Flags().StringVarP(
		&outputPath,
		"output", "o", types.CurrentWDPath(types.DefaultOutputFileName),
		`Path to the output file. The output type is inferred from the file extension. valid types:`+
			fmt.Sprintf("%v", types.OutputTypeNames()),
	)

	return cmd
}
