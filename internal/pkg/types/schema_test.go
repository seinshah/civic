package types_test

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/seinshah/cvci/internal/pkg/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	_ "embed"
)

//go:embed example.schema.yaml
var sampleSchema []byte

var re = regexp.MustCompile(`\[(?P<sliceIndex>\d+)]`)

func TestSchemaDataType(t *testing.T) {
	t.Parallel()

	var data types.Schema

	err := defaults.Set(&data)

	require.NoError(t, err)

	err = yaml.Unmarshal(sampleSchema, &data)

	require.NoError(t, err)

	validate := validator.New(validator.WithRequiredStructEnabled())

	err = validate.Struct(data)

	require.NoError(t, err)

	testCases := []struct {
		name     string
		key      string
		validate func(value any) bool
	}{
		{
			name: "template path",
			key:  "Template.Path",
			validate: func(value any) bool {
				d := value.(string)

				return d == "https://raw.githubusercontent.com/seinshah/cvci/main/examples/example.template.html"
			},
		},
		{
			name: "bio contact socials",
			key:  "Bio.Contact.Socials",
			validate: func(value any) bool {
				d := value.([]any)

				return len(d) == 2
			},
		},
		{
			name: "work experiences header",
			key:  "WorkExperiences.Header",
			validate: func(value any) bool {
				d := value.(string)

				return d == "Work Experiences2"
			},
		},
		{
			name: "education first entity details",
			key:  "Educations.Entities[0].Details",
			validate: func(value any) bool {
				d := value.([]any)

				return len(d) == 3
			},
		},
		{
			name: "skills default header",
			key:  "Skills.Header",
			validate: func(value any) bool {
				d := value.(string)

				return d == "Skills"
			},
		},
		{
			name: "second custom section details",
			key:  "CustomSections[1].Details",
			validate: func(value any) bool {
				d := value.([]any)

				return len(d) == 3
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			keys := strings.Split(tc.key, ".")

			var (
				value        any
				internalData map[string]any
				ok           bool
			)

			//nolint:musttag
			jsonData, errI := json.Marshal(data)

			require.NoError(t, errI)

			errI = json.Unmarshal(jsonData, &internalData)

			require.NoError(t, errI)

			for index, key := range keys {
				matches := re.FindStringSubmatch(key)
				keyIndex := -1

				if len(matches) == 2 {
					keyIndex, errI = strconv.Atoi(matches[1])

					require.NoError(t, errI)

					key = strings.ReplaceAll(key, matches[0], "")
				}

				if index == len(keys)-1 {
					value = internalData[key]

					if keyIndex != -1 {
						value = internalData[key].([]any)[keyIndex]
					}
				} else {
					if keyIndex != -1 {
						internalData, ok = internalData[key].([]any)[keyIndex].(map[string]any)
					} else {
						internalData, ok = internalData[key].(map[string]any)
					}

					require.True(t, ok)
				}
			}

			require.Truef(t, tc.validate(value), "valid value %v for key %s ", value, tc.key)
		})
	}
}
