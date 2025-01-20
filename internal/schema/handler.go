package schema

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/seinshah/civic/internal/pkg/loader"
	"github.com/seinshah/civic/internal/pkg/types"
)

const sampleConfigPath = "https://raw.githubusercontent.com/seinshah/civic/main/examples/example.schema.yaml"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Init creates a new configuration file using the sample values.
// Then user will be able to pick it up and modify it.
func (h *Handler) Init(
	ctx context.Context,
	outputPath string,
) error {
	if outputPath == "" {
		return types.ErrEmptyOutputPath
	}

	x := types.DetectFileType[types.SchemaType](outputPath)
	if !x.IsValid() {
		return fmt.Errorf(
			"%w: couldn't detect the file type from %s. (valid types: %v)",
			types.ErrInvalidSchemaType, outputPath, types.SchemaTypeNames(),
		)
	}

	slog.Debug("loading sample cv schema file", "path", sampleConfigPath)

	conf, err := loader.NewRemoteLoader(sampleConfigPath).Load(ctx)
	if err != nil {
		return fmt.Errorf("failed to load sample schema file: %w", err)
	}

	slog.Debug("loaded sample schema file", "writingTo", outputPath)

	if err = os.WriteFile(outputPath, conf, types.DefaultFilePermission); err != nil {
		return fmt.Errorf("failed to write the schema file: %w", err)
	}

	slog.Info("Sample schema file created successfully", "path", outputPath)

	return nil
}
