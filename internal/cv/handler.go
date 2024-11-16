package cv

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/seinshah/cvci/internal/pkg/output"
	"github.com/seinshah/cvci/internal/pkg/output/html"
	"github.com/seinshah/cvci/internal/pkg/output/pdf/chrome"
	"github.com/seinshah/cvci/internal/pkg/types"
)

type Handler struct {
	appVersion       string
	resumeConfigPath string
	outputPath       string
	outputType       types.OutputType
}

var errInvalidOutputType = errors.New("invalid file type")

func NewHandler(
	appVersion string,
	configPath string,
	outputPath string,
) (*Handler, error) {
	if outputPath == "" {
		outputPath = fmt.Sprintf(
			"%s%s%s",
			outputPath, string(os.PathSeparator),
			types.DefaultOutputFileName,
		)
	}

	outputType := types.DetectOutputType(outputPath)
	if !outputType.IsValid() {
		return nil, fmt.Errorf(
			"%w: couldn't detect the file type from %s. (valid types: %v)",
			errInvalidOutputType,
			outputPath,
			types.OutputTypeNames(),
		)
	}

	if configPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get the current working directory: %w", err)
		}

		configPath = fmt.Sprintf("%s%s%s", wd, string(os.PathSeparator), types.DefaultConfigFileName)
	}

	return &Handler{
		resumeConfigPath: configPath,
		appVersion:       appVersion,
		outputPath:       outputPath,
		outputType:       outputType,
	}, nil
}

func (h *Handler) Generate(ctx context.Context) error {
	confData, err := h.getConfigDate(ctx)
	if err != nil {
		return err
	}

	templateContent, err := h.parseTemplate(ctx, confData)
	if err != nil {
		return err
	}

	generator, err := h.getOutputGenerator(confData)
	if err != nil {
		return err
	}

	if err = output.Render(ctx, templateContent, generator, h.outputPath); err != nil {
		return fmt.Errorf("failed to render the output: %w", err)
	}

	slog.Info("Rendered the output. Your CV should be ready on " + h.outputPath)

	return nil
}

func (h *Handler) getConfigDate(ctx context.Context) (*types.Schema, error) {
	confManager, err := NewConfiguration(ctx, ConfigurationConfig{
		ConfigPath: h.resumeConfigPath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initiate configuration file: %w", err)
	}

	slog.Debug("Validating the configuration file...")

	if err = confManager.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration file: %w", err)
	}

	confData := confManager.Data()

	slog.Info("Validated and parsed the configuration file")

	return confData, nil
}

func (h *Handler) parseTemplate(ctx context.Context, confData *types.Schema) ([]byte, error) {
	slog.Debug("Creating template manager...")

	templateManager, err := NewTemplate(ctx, TemplateConfig{
		AppVersion:   h.appVersion,
		TemplatePath: confData.Template.Path,
		Customizer:   confData.Template.Customizer,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initiate template file: %w", err)
	}

	slog.Debug("Validating the template file...")

	if err = templateManager.Validate(); err != nil {
		return nil, fmt.Errorf("invalid template file: %w", err)
	}

	slog.Debug("Processing the template file...")

	templateContent, err := templateManager.Process(confData)
	if err != nil {
		return nil, fmt.Errorf("failed to process the template: %w", err)
	}

	slog.Info("Validated, parsed, and processed the template file")

	return templateContent, nil
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
		return nil, errInvalidOutputType
	}

	return generator, nil
}
