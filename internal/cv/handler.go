package cv

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/seinshah/civic/internal/pkg/output"
	"github.com/seinshah/civic/internal/pkg/output/html"
	"github.com/seinshah/civic/internal/pkg/output/pdf/chrome"
	"github.com/seinshah/civic/internal/pkg/types"
)

var ErrGenerateOutput = errors.New("failed to generate the output")

type Handler struct {
	appVersion     string
	schemaFilePath string
	schemaType     types.SchemaType
	outputPath     string
	outputType     types.OutputType
}

func NewHandler(
	appVersion string,
	schemaFilePath string,
	outputPath string,
) (*Handler, error) {
	if outputPath == "" {
		return nil, types.ErrEmptyOutputPath
	}

	outputType := types.DetectFileType[types.OutputType](outputPath)
	if !outputType.IsValid() {
		return nil, fmt.Errorf(
			"%w: couldn't detect the file type from %s. (valid types: %v)",
			types.ErrInvalidOutputType,
			outputPath,
			types.OutputTypeNames(),
		)
	}

	if schemaFilePath == "" {
		return nil, types.ErrEmptySchemaPath
	}

	schemaType := types.DetectFileType[types.SchemaType](schemaFilePath)
	if !schemaType.IsValid() {
		return nil, fmt.Errorf(
			"%w: couldn't detect the file type from %s. (valid types: %v)",
			types.ErrInvalidSchemaType,
			schemaFilePath,
			types.SchemaTypeNames(),
		)
	}

	return &Handler{
		appVersion:     appVersion,
		schemaFilePath: schemaFilePath,
		schemaType:     schemaType,
		outputPath:     outputPath,
		outputType:     outputType,
	}, nil
}

func (h *Handler) Generate(ctx context.Context) error {
	confData, err := h.parseSchemaFile(ctx)
	if err != nil {
		return err
	}

	slog.Info("Successfully processed the CV schema file")

	templateContent, err := h.parseTemplate(ctx, confData)
	if err != nil {
		return err
	}

	slog.Info("Successfully processed the template file")

	generator, err := h.getOutputGenerator(confData)
	if err != nil {
		return err
	}

	if err = output.Render(ctx, templateContent, generator, h.outputPath); err != nil {
		return errors.Join(ErrGenerateOutput, err)
	}

	slog.Info("Rendered the output. Your CV should be ready on " + h.outputPath)

	return nil
}

// getOutputGenerator returns the output generator based on the output type.
// nolint:ireturn
func (h *Handler) getOutputGenerator(confData *types.Schema) (types.OutputGenerator, error) {
	var generator types.OutputGenerator

	switch h.outputType {
	case types.OutputTypePdf:
		generator = chrome.NewHeadless(
			chrome.WithPageSize(confData.Page.Size),
			chrome.WithPageMargin(confData.Page.Margin),
		)

		slog.Debug("Rendering the PDF...")

	case types.OutputTypeHtml:
		generator = html.NewEngine()

		slog.Debug("Rendering the HTML...")

	default:
		return nil, types.ErrInvalidOutputType
	}

	return generator, nil
}
