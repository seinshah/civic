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
	"strings"

	"github.com/seinshah/cvci/internal/pkg/loader"
	"github.com/seinshah/cvci/internal/pkg/types"
	"github.com/seinshah/flattenhtml"
)

const (
	metaAttributeAppVersion        = "app-version"        // nolint: unused
	metaAttributeTemplateDirection = "template-direction" // nolint: unused
)

var (
	ErrNonParsableTemplate = errors.New("HTML template cannot be parsed")
	ErrFoundInvalidTag     = errors.New("found invalid tag in the HTML template")
	ErrMismatchAppVersion  = errors.New("template does not support the current app version")

	ErrInvalidDirective = errors.New("template file is using an unsupported directive")

	errNoCursor = errors.New("no cursor")
)

var (
	forbiddenTags = []string{"script", "iframe", "link"} // nolint: unused

	// nolint: unused
	forbiddenException = map[string]map[string][]string{
		"link": {
			"rel": []string{"stylesheet"},
		},
	}
)

func (h *Handler) parseTemplate(ctx context.Context, config *types.Schema) ([]byte, error) {
	templateLoader, err := loader.NewGeneralLoader(config.Template.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to load template file (%s): %w", config.Template.Path, err)
	}

	content, err := templateLoader.Load(ctx)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("cvci").Parse(string(content))
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

	if err = customizeTemplate(cursor, config.Template.Customizer); err != nil {
		slog.Warn("failed to register new node", "error", err, "customizer", config.Template.Customizer)
	}

	var output bytes.Buffer

	if err = nodeManager.Render(&output); err != nil {
		return nil, fmt.Errorf("failed to render the modified template: %w", err)
	}

	return output.Bytes(), nil
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
				map[string]string{"type": "text/css", "data-source": "cvci_customizer"},
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
// nolint: unused
func (t *templateValidator) ValidateForbiddenTags() error {
	for _, tag := range forbiddenTags {
		tags := t.cursor.SelectNodes(tag)

		exceptions, ok := forbiddenException[tag]

		// TODO: Add predicates to flattenhtml for filter functionality
		// TODO: Add error return type to (*Nodes).Each method
		if ok {
			var invalidTags []string

			for attribute, exceptionValues := range exceptions {
				tags.Each(func(node *flattenhtml.Node) {
					nodeAttrVal, ok := node.Attribute(attribute)
					if !ok || !slices.Contains(exceptionValues, nodeAttrVal) {
						invalidTags = append(
							invalidTags,
							fmt.Sprintf("%s(%v)", tag, node.Attributes()),
						)
					}
				})

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
// nolint: unused
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

	appMajor := strings.TrimPrefix(strings.Split(t.appVersion, ".")[0], "v")
	templateAppMajor := strings.TrimPrefix(strings.Split(tplAppVersion, ".")[0], "v")

	if appMajor != templateAppMajor {
		return fmt.Errorf(
			"app version mismatch (%s != %s): %w",
			appMajor, templateAppMajor, ErrMismatchAppVersion,
		)
	}

	return nil
}
