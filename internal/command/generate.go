package command

import (
	"github.com/seinshah/cvci/internal/cv"
	"github.com/spf13/cobra"
)

func (c *Command) getGenerateCommands() *cobra.Command {
	var (
		configPath string
		outputPath string
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate the resume or cv",
		Long:  `Generate the resume or cv based on the configuration file and the provided version.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			handler, err := cv.NewHandler(cmd.Version, configPath, outputPath)
			if err != nil {
				return err
			}

			return handler.Generate(cmd.Context())
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
