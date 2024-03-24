package config

import (
	"context"
	"fmt"
	"os"

	"github.com/seinshah/cvci/internal/pkg/loader"
	"github.com/seinshah/cvci/internal/pkg/types"
)

const sampleConfigPath = "https://raw.githubusercontent.com/seinshah/cvci/main/assets/sample_config.yaml"

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
		outputPath = "./" + types.DefaultConfigFileName
	}

	conf, err := loader.NewRemoteLoader(sampleConfigPath).Load(ctx)
	if err != nil {
		return fmt.Errorf("failed to load sample configuration file: %w", err)
	}

	if err = os.WriteFile(outputPath, conf, types.DefaultFilePermission); err != nil {
		return fmt.Errorf("failed to write the configuration file: %w", err)
	}

	return nil
}
