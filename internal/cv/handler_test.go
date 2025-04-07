package cv_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/seinshah/civic/internal/cv"
	"github.com/seinshah/civic/internal/pkg/loader"
	"github.com/seinshah/civic/internal/pkg/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// getMapKeyValue accepts a nested interface map and . separated key string
// and returns the value of the leaf key from the map.
func getMapKeyValue(t *testing.T, m map[string]any, key string) any {
	t.Helper()

	if m == nil {
		return nil
	}

	keys := strings.SplitN(key, ".", 2)

	val, ok := m[keys[0]]
	if !ok {
		return nil
	}

	if len(keys) == 1 {
		return val
	}

	valMap, ok := val.(map[string]any)
	if ok {
		return getMapKeyValue(t, valMap, keys[1])
	}

	t.Fatalf("key %s requested on map %s, which is not a map", keys[1], keys[0])

	return nil
}

func getSchemaPath(t *testing.T, schemaContent map[string]any, templateContent string) string {
	t.Helper()

	var (
		err            error
		createTemplate bool
		templateFile   *os.File
		schemaFile     *os.File
	)

	// replace template.path with the path to the created template file
	// if it's provided value is "<<template_path>>"
	if pathVal := getMapKeyValue(t, schemaContent, "template.path"); pathVal != nil && pathVal.(string) == "<<template_path>>" {
		createTemplate = true
	}

	if createTemplate {
		templateFile, err = os.CreateTemp(os.TempDir(), "app-test-template*.html")

		require.NoError(t, err)

		_, err = templateFile.WriteString(templateContent)

		require.NoError(t, err)

		schemaContent["template"].(map[string]any)["path"] = templateFile.Name()
	}

	yamlSchema, err := yaml.Marshal(schemaContent)

	require.NoError(t, err)

	schemaFile, err = os.CreateTemp(os.TempDir(), "app-test-schema*.yaml")

	require.NoError(t, err)

	_, err = schemaFile.Write(yamlSchema)

	require.NoError(t, err)

	t.Cleanup(func() {
		err = os.Remove(schemaFile.Name())

		require.NoError(t, err)

		if createTemplate {
			err = os.Remove(templateFile.Name())

			require.NoError(t, err)
		}
	})

	return schemaFile.Name()
}

func TestNewHandler(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		outputPath     string
		schemaFilePath string
		hasError       bool
		err            error
	}{
		{
			name:     "empty output path",
			hasError: true,
			err:      types.ErrEmptyOutputPath,
		},
		{
			name:       "invalid output file type",
			outputPath: "output.jpeg",
			hasError:   true,
			err:        types.ErrInvalidOutputType,
		},
		{
			name:       "empty schema file path",
			outputPath: "output.pdf",
			hasError:   true,
			err:        types.ErrEmptySchemaPath,
		},
		{
			name:           "invalid schema file type",
			outputPath:     "output.pdf",
			schemaFilePath: "schema.json",
			hasError:       true,
			err:            types.ErrInvalidSchemaType,
		},
		{
			name:           "valid handler",
			outputPath:     "output.pdf",
			schemaFilePath: "schema.yaml",
		},
		{
			name:           "valid handler with different output type",
			outputPath:     "output.html",
			schemaFilePath: "schema.yaml",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h, err := cv.NewHandler("v0.1.0", tc.schemaFilePath, tc.outputPath)
			if tc.hasError {
				require.Error(t, err)
				require.Nil(t, h)

				if tc.err != nil {
					require.ErrorIs(t, err, tc.err)
				}

				return
			}

			require.NoError(t, err)
			require.NotNil(t, h)
		})
	}
}

func TestHandler_Generate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		schemaFilePath  func(*testing.T) string
		outputFile      string
		outputExtension string
		ctx             func(t *testing.T) (context.Context, context.CancelFunc)
		validateOutput  func(t *testing.T, outputFile string)
		hasError        bool
		err             error
	}{
		{
			name: "non-existent schema file",
			schemaFilePath: func(t *testing.T) string {
				t.Helper()

				return "non-existent.yaml"
			},
			hasError: true,
			err:      loader.ErrInvalidLocalPath,
		},
		{
			name: "invalid schema file",
			schemaFilePath: func(t *testing.T) string {
				t.Helper()

				return getSchemaPath(t, nil, "")
			},
			hasError: true,
			err:      cv.ErrInvalidSchemaFormat,
		},
		{
			name: "schema file with invalid value",
			schemaFilePath: func(t *testing.T) string {
				t.Helper()

				return getSchemaPath(
					t,
					map[string]any{
						"template": map[string]any{
							"path": "<<template_path>>",
						},
						"bio": map[string]any{
							"name":  "John Doe",
							"title": "S", // min 2 characters
						},
					},
					"",
				)
			},
			hasError: true,
			err:      cv.ErrInvalidSchemaFormat,
		},
		{
			name: "non existing template file",
			schemaFilePath: func(t *testing.T) string {
				t.Helper()

				return getSchemaPath(
					t,
					map[string]any{
						"template": map[string]any{
							"path": os.TempDir() + "/non-existent.html",
						},
						"bio": map[string]any{
							"name":  "John Doe",
							"title": "Software Engineer",
						},
					},
					"",
				)
			},
			hasError: true,
			err:      loader.ErrInvalidLocalPath,
		},
		{
			name: "with invalid template directive",
			schemaFilePath: func(t *testing.T) string {
				t.Helper()

				return getSchemaPath(
					t,
					map[string]any{
						"template": map[string]any{
							"path": "<<template_path>>",
						},
						"bio": map[string]any{
							"name":  "John Doe",
							"title": "Software Engineer",
						},
					},
					"{{.Raw.SomeInvalidDirective}}",
				)
			},
			hasError: true,
			err:      cv.ErrInvalidDirective,
		},
		{
			name: "forbidden html tag link",
			schemaFilePath: func(t *testing.T) string {
				t.Helper()

				return getSchemaPath(
					t,
					map[string]any{
						"template": map[string]any{
							"path": "<<template_path>>",
						},
						"bio": map[string]any{
							"name":  "John Doe",
							"title": "Software Engineer",
						},
					},
					`
						<meta name="app-version" content="v0" />
						<link rel="disallowed" href="style.css" />
					`,
				)
			},
			hasError: true,
			err:      cv.ErrFoundInvalidTag,
		},
		{
			name: "forbidden html tag script",
			schemaFilePath: func(t *testing.T) string {
				t.Helper()

				return getSchemaPath(
					t,
					map[string]any{
						"template": map[string]any{
							"path": "<<template_path>>",
						},
						"bio": map[string]any{
							"name":  "John Doe",
							"title": "Software Engineer",
						},
					},
					`
						<meta name="app-version" content="v0" />
						<script src="script.js"></script>
					`,
				)
			},
			hasError: true,
			err:      cv.ErrFoundInvalidTag,
		},
		{
			name: "forbidden html tag iframe",
			schemaFilePath: func(t *testing.T) string {
				t.Helper()

				return getSchemaPath(
					t,
					map[string]any{
						"template": map[string]any{
							"path": "<<template_path>>",
						},
						"bio": map[string]any{
							"name":  "John Doe",
							"title": "Software Engineer",
						},
					},
					`
						<meta name="app-version" content="v0" />
						<iframe src="https://example.com"></iframe>
					`,
				)
			},
			hasError: true,
			err:      cv.ErrFoundInvalidTag,
		},
		{
			name: "invalid app version",
			schemaFilePath: func(t *testing.T) string {
				t.Helper()

				return getSchemaPath(
					t,
					map[string]any{
						"template": map[string]any{
							"path": "<<template_path>>",
						},
						"bio": map[string]any{
							"name":  "John Doe",
							"title": "Software Engineer",
						},
					},
					`
						<meta name="app-version" content="v1" />
						<link rel="stylesheet" href="style.css" />
					`,
				)
			},
			hasError: true,
			err:      cv.ErrMismatchAppVersion,
		},
		{
			name:       "invalid output path",
			outputFile: "/invalid/output/path/output.pdf",
			schemaFilePath: func(t *testing.T) string {
				t.Helper()

				return getSchemaPath(
					t,
					map[string]any{
						"template": map[string]any{
							"path": "<<template_path>>",
						},
						"bio": map[string]any{
							"name":  "John Doe",
							"title": "Software Engineer",
						},
					},
					`
						<meta name="app-version" content="v0" />
						<link rel="stylesheet" href="style.css" />
					`,
				)
			},
			hasError: true,
			err:      cv.ErrGenerateOutput,
		},
		{
			name: "context error",
			schemaFilePath: func(t *testing.T) string {
				t.Helper()

				return getSchemaPath(
					t,
					map[string]any{
						"template": map[string]any{
							"path": "https://example.com/template.html",
						},
						"bio": map[string]any{
							"name":  "John Doe",
							"title": "Software Engineer",
						},
					},
					`
						<html>
							<head>
								<meta name="app-version" content="v0" />
								<link rel="stylesheet" href="style.css" />
							</head>
							<body>
								<h1>{{.Raw.Bio.Name}}</h1>
								<h2>{{.Raw.Bio.Title}}</h2>
							</body>
						</html>
					`,
				)
			},
			ctx: func(t *testing.T) (context.Context, context.CancelFunc) {
				t.Helper()

				return context.WithTimeout(context.Background(), 0)
			},
			hasError: true,
			err:      context.DeadlineExceeded,
		},
		{
			name:            "valid html output",
			outputExtension: "html",
			schemaFilePath: func(t *testing.T) string {
				t.Helper()

				return getSchemaPath(
					t,
					map[string]any{
						"template": map[string]any{
							"path": "<<template_path>>",
						},
						"bio": map[string]any{
							"name":  "John Doe",
							"title": "Software Engineer",
						},
					},
					`
						<html>
							<head>
								<meta name="app-version" content="v0" />
								<link rel="stylesheet" href="style.css" />
							</head>
							<body>
								<h1>{{.Raw.Bio.Name}}</h1>
								<h2>{{.Raw.Bio.Title}}</h2>
								{{with .Raw.WorkExperiences}}<h3>Work Experiences</h3>{{end}}
								{{with .Raw.Educations}}<h3>Educations</h3>{{end}}
								{{with .Raw.Certificates}}<h3>Certificates</h3>{{end}}
								{{with .Raw.Publications}}<h3>Publications</h3>{{end}}
								{{with .Raw.Skills}}<h3>Skills</h3>{{end}}
								{{with .Raw.CustomSections}}<h3>CustomSections</h3>{{end}}
							</body>
						</html>
					`,
				)
			},
			validateOutput: func(t *testing.T, outputFile string) {
				t.Helper()

				data, err := os.ReadFile(outputFile)

				require.NoError(t, err)

				require.Containsf(t, string(data), "<h1>John Doe</h1>", "Name should be present in the output")
				require.Containsf(t, string(data), "<h2>Software Engineer</h2>", "Title should be present in the output")

				require.NotContainsf(t, string(data), "<h3>Work Experiences</h3>", "Work Experiences should not be present in the output")
				require.NotContainsf(t, string(data), "<h3>Educations</h3>", "Educations should not be present in the output")
				require.NotContainsf(t, string(data), "<h3>Certificates</h3>", "Certificates should not be present in the output")
				require.NotContainsf(t, string(data), "<h3>Publications</h3>", "Publications should not be present in the output")
				require.NotContainsf(t, string(data), "<h3>Skills</h3>", "Skills should not be present in the output")
				require.NotContainsf(t, string(data), "<h3>CustomSections</h3>", "CustomSections should not be present in the output")
			},
		},
		{
			name:            "valid pdf output",
			outputExtension: "pdf",
			schemaFilePath: func(t *testing.T) string {
				t.Helper()

				return getSchemaPath(
					t,
					map[string]any{
						"template": map[string]any{
							"path": "<<template_path>>",
						},
						"bio": map[string]any{
							"name":  "John Doe",
							"title": "Software Engineer",
						},
					},
					`
						<html>
							<head>
								<meta name="app-version" content="v0" />
								<link rel="stylesheet" href="style.css" />
							</head>
							<body>
								<h1>{{.Raw.Bio.Name}}</h1>
								<h2>{{.Raw.Bio.Title}}</h2>
							</body>
						</html>
					`,
				)
			},
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.outputFile == "" {
				ext := "html"
				if tc.outputExtension != "" {
					ext = tc.outputExtension
				}

				tc.outputFile = fmt.Sprintf("%s/app-generate-test-output-%d.%s", os.TempDir(), i, ext)
			}

			h, err := cv.NewHandler("v0.1.0", tc.schemaFilePath(t), tc.outputFile)
			require.NoError(t, err)
			require.NotNil(t, h)

			ctx := context.Background()

			if tc.ctx != nil {
				var cancel context.CancelFunc

				ctx, cancel = tc.ctx(t)

				defer cancel()
			}

			err = h.Generate(ctx)
			if tc.hasError {
				require.Error(t, err)

				if tc.err != nil {
					require.ErrorIs(t, err, tc.err)
				}

				return
			}

			require.NoError(t, err)
			require.FileExists(t, tc.outputFile)

			if tc.validateOutput != nil {
				tc.validateOutput(t, tc.outputFile)
			}
		})
	}
}
