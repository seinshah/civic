package cvhandler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/seinshah/cvci/internal/pkg/loader"
	"github.com/seinshah/cvci/internal/pkg/types"
	"github.com/seinshah/flattenhtml"
)

const (
	metaAttributeAppVersion        = "app-version"
	metaAttributeTemplateDirection = "template-direction"
	metaAttributeTemplateLanguage  = "template-language"
)

type TemplateConfig struct {
	// AppVersion is the current version of the application.
	AppVersion string `validate:"required,semver"`

	// TemplatePath is the path to the HTML template.
	TemplatePath string `validate:"required"`

	// Customizer is the manual customization that will be added to
	// the template.
	Customizer types.Customizer
}

type Template struct {
	content     []byte
	cursor      *flattenhtml.Cursor
	nodeManager *flattenhtml.NodeManager
	config      TemplateConfig
}

var (
	ErrNonParsableTemplate = errors.New("HTML template cannot be parsed")
	ErrFoundInvalidTag     = errors.New("found invalid tag in the HTML template")
	ErrMismatchAppVersion  = errors.New("template does not support the current app version")
)

var forbiddenTags = []string{"script", "iframe", "link"}

func NewTemplate(ctx context.Context, config TemplateConfig) (*Template, error) {
	templateLoader, err := loader.NewGeneralLoader(config.TemplatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template file (%s): %w", config.TemplatePath, err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("invalid template config: %w", err)
	}

	content, err := templateLoader.Load(ctx)
	if err != nil {
		return nil, err
	}

	nodeManager, err := flattenhtml.NewNodeManagerFromReader(bytes.NewReader(content))
	if err != nil {
		return nil, errors.Join(ErrNonParsableTemplate, err)
	}

	multiCursor, err := nodeManager.Parse(flattenhtml.NewTagFlattener())
	if err != nil {
		return nil, errors.Join(ErrNonParsableTemplate, err)
	}

	cursor := multiCursor.First()
	for cursor != nil {
		return nil, fmt.Errorf("no cursor: %w", ErrNonParsableTemplate)
	}

	return &Template{
		content:     content,
		cursor:      cursor,
		nodeManager: nodeManager,
		config:      config,
	}, nil
}

func (t *Template) Validate() error {
	validators := []func() error{
		t.validateForbiddenTags,
		t.validateAppVersion,
	}

	for _, tplValidator := range validators {
		if err := tplValidator(); err != nil {
			return err
		}
	}

	return t.processModifications()
}

func (t *Template) Process(config *types.Schema) ([]byte, error) {
	tpl, err := template.New("cvci").Parse(string(t.content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse the template: %w", err)
	}

	var processedTempl bytes.Buffer

	if err = tpl.Execute(&processedTempl, config); err != nil {
		return nil, fmt.Errorf("failed to execute the template: %w", err)
	}

	return processedTempl.Bytes(), nil
}

func (t *Template) processModifications() error {
	var modified bool

	if t.config.Customizer.Style != "" {
		if node := t.cursor.SelectNodes("head").First(); node != nil {
			styleTag := node.AppendChild(
				flattenhtml.NodeTypeElement,
				"style",
				map[string]string{"type": "text/css", "data-source": "cvci_customizer"},
			)

			styleTag.AppendChild(flattenhtml.NodeTypeText, t.config.Customizer.Style, nil)

			if err := t.cursor.RegisterNewNode(styleTag); err != nil {
				slog.Warn("failed to register new node", "error", err)
			}

			modified = true
		}
	}

	if modified {
		var outBuffer bytes.Buffer

		if err := t.nodeManager.Render(&outBuffer); err != nil {
			return fmt.Errorf("failed to render the modified template: %w", err)
		}

		t.content = outBuffer.Bytes()
	}

	return nil
}

func (t *Template) validateForbiddenTags() error {
	for _, tag := range forbiddenTags {
		if t.cursor.SelectNodes(tag).Len() > 0 {
			return fmt.Errorf("%s: %w", tag, ErrFoundInvalidTag)
		}
	}

	return nil
}

func (t *Template) validateAppVersion() error {
	metaTag := t.cursor.SelectNodes("meta").
		Filter(
			flattenhtml.WithAttributeValueAs("name", metaAttributeAppVersion),
		).
		First()

	if metaTag == nil {
		return fmt.Errorf("missing meta tag: %w", ErrMismatchAppVersion)
	}

	tplAppVersion, _ := metaTag.Attribute("content")
	if tplAppVersion == "" {
		return fmt.Errorf("empty %s: %w", metaAttributeAppVersion, ErrMismatchAppVersion)
	}

	appMajor := strings.TrimPrefix(strings.Split(t.config.AppVersion, ".")[0], "v")
	templateAppMajor := strings.TrimPrefix(strings.Split(tplAppVersion, ".")[0], "v")

	if appMajor != templateAppMajor {
		return fmt.Errorf(
			"app version mismatch (%s != %s): %w",
			appMajor, templateAppMajor, ErrMismatchAppVersion,
		)
	}

	return nil
}
