package cvhandler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/seinshah/cvci/internal/pkg/loader"
	"github.com/seinshah/flattenhtml"
)

const (
	metaAttributeAppVersion        = "app-version"
	metaAttributeTemplateDirection = "template-direction"
	metaAttributeTemplateLanguage  = "template-language"
)

type Config struct {
	// AppVersion is the current version of the application.
	AppVersion string `validate:"required,semver"`

	// Loader is the loader that is ready to load the HTML template.
	Loader loader.Loader `validate:"required"`
}

type Engine struct {
	content []byte
	cursor  *flattenhtml.Cursor
	config  Config
}

var (
	ErrNonParsableTemplate = errors.New("HTML template cannot be parsed")
	ErrFoundInvalidTag     = errors.New("found invalid tag in the HTML template")
	ErrMismatchAppVersion  = errors.New("template does not support the current app version")
)

var forbiddenTags = []string{"script", "iframe", "link"}

func NewEngine(ctx context.Context, config Config) (*Engine, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("invalid template config: %w", err)
	}

	content, err := config.Loader.Load(ctx)
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

	return &Engine{
		content: content,
		cursor:  cursor,
		config:  config,
	}, nil
}

func (e *Engine) Validate() error {
	validators := []func() error{
		e.validateForbiddenTags,
		e.validateAppVersion,
	}

	for _, tplValidator := range validators {
		if err := tplValidator(); err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) Process() error {
	return nil
}

func (e *Engine) validateForbiddenTags() error {
	for _, tag := range forbiddenTags {
		if e.cursor.SelectNodes(tag).Len() > 0 {
			return fmt.Errorf("%s: %w", tag, ErrFoundInvalidTag)
		}
	}

	return nil
}

func (e *Engine) validateAppVersion() error {
	metaTag := e.cursor.SelectNodes("meta").
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

	appMajor := strings.TrimPrefix(strings.Split(e.config.AppVersion, ".")[0], "v")
	templateAppMajor := strings.TrimPrefix(strings.Split(tplAppVersion, ".")[0], "v")

	if appMajor != templateAppMajor {
		return fmt.Errorf(
			"app version mismatch (%s != %s): %w",
			appMajor, templateAppMajor, ErrMismatchAppVersion,
		)
	}

	return nil
}
