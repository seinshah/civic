package schema

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/seinshah/cvci/internal/pkg/loader"
	"github.com/seinshah/cvci/internal/pkg/types"
)

const sampleConfigPath = "https://raw.githubusercontent.com/seinshah/cvci/main/examples/example.schema.yaml"

var ErrEmptyOutputPath = errors.New("output path is empty")

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
		return ErrEmptyOutputPath
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
