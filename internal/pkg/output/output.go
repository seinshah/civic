package output

import (
	"context"
	"fmt"
	"os"

	"github.com/seinshah/cvci/internal/pkg/types"
)

func Render(
	ctx context.Context,
	content []byte,
	engine types.OutputGenerator,
	outputPath string,
) error {
	output, err := engine.Generate(ctx, content)
	if err != nil {
		return err
	}

	if err = os.WriteFile(outputPath, output, types.DefaultFilePermission); err != nil {
		return fmt.Errorf("failed to write PDF to file: %w", err)
	}

	return nil
}
