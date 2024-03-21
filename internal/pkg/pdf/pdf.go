package pdf

import (
	"context"
	"fmt"
	"os"
)

type Generator interface {
	Generate(ctx context.Context, content []byte) ([]byte, error)
}

const outputPermission = 0o600

func Render(
	ctx context.Context,
	content []byte,
	engine Generator,
	outputPath string,
) error {
	output, err := engine.Generate(ctx, content)
	if err != nil {
		return err
	}

	if err = os.WriteFile(outputPath, output, outputPermission); err != nil {
		return fmt.Errorf("failed to write PDF to file: %w", err)
	}

	return nil
}
