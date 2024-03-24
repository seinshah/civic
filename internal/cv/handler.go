package cv

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/seinshah/cvci/internal/pkg/pdf"
	"github.com/seinshah/cvci/internal/pkg/pdf/chrome"
	"github.com/seinshah/cvci/internal/pkg/types"
)

type Handler struct {
	appVersion       string
	resumeConfigPath string
	outputPath       string
}

func NewHandler(
	appVersion string,
	configPath string,
	outputPath string,
) (*Handler, error) {
	if !strings.HasSuffix(outputPath, ".pdf") {
		outputPath = fmt.Sprintf("%s%s%s.pdf", outputPath, string(os.PathSeparator), types.DefaultOutputFileName)
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
	}, nil
}

func (g *Handler) Generate(ctx context.Context) error {
	confManager, err := NewConfiguration(ctx, ConfigurationConfig{
		ConfigPath: g.resumeConfigPath,
	})
	if err != nil {
		return fmt.Errorf("failed to initiate configuration file: %w", err)
	}

	slog.Debug("Validating the configuration file...")

	if err = confManager.Validate(); err != nil {
		return fmt.Errorf("invalid configuration file: %w", err)
	}

	confData := confManager.Data()

	slog.Info("Validated and parsed the configuration file")

	slog.Debug("Creating template manager...")

	templateManager, err := NewTemplate(ctx, TemplateConfig{
		AppVersion:   g.appVersion,
		TemplatePath: confData.Template.Path,
		Customizer:   confData.Template.Customizer,
	})
	if err != nil {
		return fmt.Errorf("failed to initiate template file: %w", err)
	}

	slog.Debug("Validating the template file...")

	if err = templateManager.Validate(); err != nil {
		return fmt.Errorf("invalid template file: %w", err)
	}

	slog.Debug("Processing the template file...")

	templateContent, err := templateManager.Process(confData)
	if err != nil {
		return fmt.Errorf("failed to process the template: %w", err)
	}

	slog.Info("Validated, parsed, and processed the template file")

	pdfEngine := chrome.NewHeadless(
		chrome.WithPageSize(confData.Page.Size),
		chrome.WithPageMargin(confData.Page.Margin),
	)

	slog.Debug("Rendering the PDF...")

	if err = pdf.Render(ctx, templateContent, pdfEngine, g.outputPath); err != nil {
		return fmt.Errorf("failed to render the PDF: %w", err)
	}

	slog.Info("Rendered the output. Check " + g.outputPath)

	return nil
}
