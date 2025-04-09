package cv

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"reflect"
	"slices"

	"github.com/Masterminds/sprig/v3"
	"github.com/seinshah/civic/internal/pkg/loader"
	"github.com/seinshah/civic/internal/pkg/types"
	"github.com/seinshah/civic/internal/pkg/version"
	"github.com/seinshah/flattenhtml"
)

const (
	metaAttributeAppVersion = "app-version"

	//nolint:unused
	metaAttributeTemplateDirection = "template-direction"
)

var (
	ErrTemplateNotProvided = errors.New("template path or name is required")
	ErrNonParsableTemplate = errors.New("HTML template cannot be parsed")
	ErrFoundInvalidTag     = errors.New("found invalid tag in the HTML template")
	ErrMismatchAppVersion  = errors.New("template does not support the current app version")

	ErrInvalidDirective = errors.New("template file is using an unsupported directive")

	errNoCursor = errors.New("no cursor")
)

var (
	//nolint: gochecknoglobals
	forbiddenTags = []string{"script", "iframe", "link"}

	//nolint: gochecknoglobals
	forbiddenException = map[string]map[string][]string{
		"link": {
			"rel": []string{"stylesheet"},
		},
	}
)

func (h *Handler) parseTemplate(ctx context.Context, config types.TemplateData) ([]byte, error) {
	content, err := h.getTemplateContent(ctx, config)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New(types.DefaultAppName).Funcs(sprig.FuncMap()).Parse(string(content))
	if err != nil {
		slog.Debug("", "template", string(content))

		return nil, fmt.Errorf("failed to parse the template: %w", err)
	}

	var processedTemplate bytes.Buffer

	if err = tpl.Execute(&processedTemplate, config); err != nil {
		slog.Debug("", "template", string(content), "data", config)

		return nil, errors.Join(ErrInvalidDirective, err)
	}

	nodeManager, cursor, err := initiateFlattener(&processedTemplate)
	if err != nil {
		return nil, errors.Join(ErrNonParsableTemplate, err)
	}

	if err = runTemplateValidations(cursor, h.appVersion); err != nil {
		return nil, err
	}

	if err = customizeTemplate(cursor, config.Raw.Template.Customizer); err != nil {
		slog.Warn("failed to register new node", "error", err, "customizer", config.Raw.Template.Customizer)
	}

	var output bytes.Buffer

	if err = nodeManager.Render(&output); err != nil {
		return nil, fmt.Errorf("failed to render the modified template: %w", err)
	}

	return output.Bytes(), nil
}

func (h *Handler) getTemplateContent(ctx context.Context, config types.TemplateData) ([]byte, error) {
	if config.Raw.Template.Path == "" && config.Raw.Template.Name == "" {
		return nil, ErrTemplateNotProvided
	}

	var templatePath string

	if config.Raw.Template.Path != "" {
		templatePath = config.Raw.Template.Path
	} else if config.Raw.Template.Name != "" {
		appV, err := version.Parse(h.appVersion)
		if err != nil {
			return nil, fmt.Errorf("invalid app version: %w", err)
		}

		templatePath = fmt.Sprintf(
			"%s/%s/v%d/template.html",
			types.TemplateRegistryPath,
			config.Raw.Template.Name,
			appV.Major(),
		)
	}

	templateLoader, err := loader.NewGeneralLoader(templatePath)
	if err != nil {
		if !errors.Is(err, loader.ErrInvalidPath) {
			return nil, fmt.Errorf("failed to load template file (%s): %w", config.Raw.Template.Path, err)
		}
	}

	return templateLoader.Load(ctx)
}

func initiateFlattener(data io.Reader) (*flattenhtml.NodeManager, *flattenhtml.Cursor, error) {
	nodeManager, err := flattenhtml.NewNodeManagerFromReader(data)
	if err != nil {
		return nil, nil, err
	}

	multiCursor, err := nodeManager.Parse(flattenhtml.NewTagFlattener())
	if err != nil {
		return nil, nil, err
	}

	cursor := multiCursor.First()
	for cursor == nil {
		return nil, nil, errNoCursor
	}

	return nodeManager, cursor, nil
}

func customizeTemplate(htmlCursor *flattenhtml.Cursor, customizer types.Customizer) error {
	if customizer.Style != "" {
		if node := htmlCursor.SelectNodes("head").First(); node != nil {
			styleTag := node.AppendChild(
				flattenhtml.NodeTypeElement,
				"style",
				map[string]string{"type": "text/css", "data-source": types.DefaultAppName + "_customizer"},
			)

			styleTag.AppendChild(flattenhtml.NodeTypeText, customizer.Style, nil)

			if err := htmlCursor.RegisterNewNode(styleTag); err != nil {
				return err
			}
		}
	}

	return nil
}

func runTemplateValidations(htmlCursor *flattenhtml.Cursor, appVersion string) error {
	v := &templateValidator{
		cursor:     htmlCursor,
		appVersion: appVersion,
	}

	vv := reflect.ValueOf(v)

	for i := range vv.NumMethod() {
		validatorMethod := vv.Method(i)

		if validatorMethod.Kind() != reflect.Func {
			continue
		} else if validatorMethod.Type().NumIn() != 0 || validatorMethod.Type().NumOut() != 1 {
			continue
		}

		mOut := validatorMethod.Call(nil)

		if len(mOut) == 0 {
			continue
		}

		if err, ok := mOut[0].Interface().(error); ok {
			return err
		}
	}

	return nil
}

// templateValidator is an internal type wrapper to define template's tag validators.
// Any exposed method defined on this type with no input arguments and an error return type,
// will be executed during runTemplateValidations.
type templateValidator struct {
	cursor     *flattenhtml.Cursor
	appVersion string
}

// ValidateForbiddenTags checks if the provided template includes any forbidden tag listed
// in forbiddenTags. It considers forbiddenException and ignore scenarios depicted in the map.
func (t *templateValidator) ValidateForbiddenTags() error {
	for _, tag := range forbiddenTags {
		tags := t.cursor.SelectNodes(tag)

		exceptions, ok := forbiddenException[tag]

		// TODO: Add predicates to flattenhtml for filter functionality
		// TODO: Add error return type to (*Nodes).Each method
		if ok {
			var invalidTags []string

			for attribute, exceptionValues := range exceptions {
				tags.Each(
					func(node *flattenhtml.Node) {
						nodeAttrVal, ok := node.Attribute(attribute)
						if !ok || !slices.Contains(exceptionValues, nodeAttrVal) {
							invalidTags = append(
								invalidTags,
								fmt.Sprintf("%s(%v)", tag, node.Attributes()),
							)
						}
					},
				)

				if len(invalidTags) > 0 {
					return fmt.Errorf("%w: %v", ErrFoundInvalidTag, invalidTags)
				}
			}
		} else if tags.Len() > 0 {
			return fmt.Errorf("%s: %w", tag, ErrFoundInvalidTag)
		}
	}

	return nil
}

// ValidateAppVersion checks if the provided template supports the current app version.
// It does so by comparing the major version of the app with the major version of the template.
func (t *templateValidator) ValidateAppVersion() error {
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

	templateV, err := version.Parse(tplAppVersion)
	if err != nil {
		return fmt.Errorf("invalid %s: %w", metaAttributeAppVersion, ErrMismatchAppVersion)
	}

	appV, err := version.Parse(t.appVersion)
	if err != nil {
		return fmt.Errorf("invalid app version: %w", err)
	}

	if appV.Major() != templateV.Major() {
		return fmt.Errorf("expected v%d template version: %w", appV.Major(), ErrMismatchAppVersion)
	}

	if !appV.Equal(templateV) {
		slog.Warn(
			"template's minor/patch version does not match the app version and in rare cases might lead to unexpected behavior",
			"template", templateV, "app", appV,
		)
	}

	return nil
}
