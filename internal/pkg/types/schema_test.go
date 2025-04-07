package types_test

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/seinshah/civic/internal/pkg/types"
	"github.com/stretchr/testify/require"

	_ "embed"
)

//go:embed example.schema.yaml
var sampleSchema []byte

var re = regexp.MustCompile(`\[(?P<sliceIndex>\d+)]`)

func TestNewSchema(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		content      string
		hasInitError bool
		// all failed validation where key is the full field name and value is the failed rule
		failedValidationFields map[string]string
	}{
		{
			name:         "invalid content format",
			content:      "invalid content",
			hasInitError: true,
		},
		{
			name:    "empty content",
			content: "",
			failedValidationFields: map[string]string{
				"Schema.Template": "required",
				"Schema.Bio":      "required",
			},
		},
		{
			name:    "empty template path and name",
			content: `template: {customizer: {style: "style.css"}}`,
			failedValidationFields: map[string]string{
				"Schema.Template.Path": "required_without",
				"Schema.Template.Name": "required_without",
				"Schema.Bio":           "required",
			},
		},
		{
			name:    "with template path",
			content: `template: {path: "path", customizer: {style: "style.css"}}`,
			failedValidationFields: map[string]string{
				"Schema.Bio": "required",
			},
		},
		{
			name:    "with template name",
			content: `template: {name: "name", customizer: {style: "style.css"}}`,
			failedValidationFields: map[string]string{
				"Schema.Bio": "required",
			},
		},
		{
			name:    "empty bio name and title",
			content: `bio: {about: "a"}`,
			failedValidationFields: map[string]string{
				"Schema.Template":  "required",
				"Schema.Bio.Name":  "required",
				"Schema.Bio.Title": "required",
			},
		},
		{
			name:    "invalid bio name and title length",
			content: `bio: {name: "a", title: "b"}`,
			failedValidationFields: map[string]string{
				"Schema.Template":  "required",
				"Schema.Bio.Name":  "min",
				"Schema.Bio.Title": "min",
			},
		},
		{
			name: "minimal valid",
			content: `template: {path: "path"}
bio: {name: "ho", title: "title"}`,
		},
		{
			name: "invalid bio contact email and customData",
			content: `template: {path: "path"}
bio: {name: "ho", title: "title", contact: {email: "x"}, customData: [{label: "x"}]}`,
			failedValidationFields: map[string]string{
				"Schema.Bio.Contact.Email":       "email",
				"Schema.Bio.CustomData[0].Value": "required",
			},
		},
		{
			name: "empty optional blocks entities",
			content: `template: {path: "path"}
bio: {name: "ho", title: "title"}
workExperiences: {entities: []}
educations: {entities: []}
certificates: {entities: []}
publications: {entities: []}
skills: {entities: []}
projects: {entities: []}
customSections: [{}]`,
			failedValidationFields: map[string]string{
				"Schema.WorkExperiences.Entities":  "min",
				"Schema.Educations.Entities":       "min",
				"Schema.Certificates.Entities":     "min",
				"Schema.Publications.Entities":     "min",
				"Schema.Skills.Entities":           "min",
				"Schema.Projects.Entities":         "min",
				"Schema.CustomSections[0].Header":  "required",
				"Schema.CustomSections[0].Details": "required",
			},
		},
		{
			name: "invalid work experiences entity",
			content: `template: {path: "path"}
bio: {name: "ho", title: "title"}
workExperiences: {entities: [{details: ["a"], technologies: [""]}]}`,
			failedValidationFields: map[string]string{
				"Schema.WorkExperiences.Entities[0].Title":           "required",
				"Schema.WorkExperiences.Entities[0].Company":         "required",
				"Schema.WorkExperiences.Entities[0].StartDate":       "required",
				"Schema.WorkExperiences.Entities[0].Details[0]":      "min",
				"Schema.WorkExperiences.Entities[0].Technologies[0]": "min",
			},
		},
		{
			name: "invalid educations entity",
			content: `template: {path: "path"}
bio: {name: "ho", title: "title"}
educations: {entities: [{details: ["a"], technologies: [""]}]}`,
			failedValidationFields: map[string]string{
				"Schema.Educations.Entities[0].Degree":          "required",
				"Schema.Educations.Entities[0].Field":           "required",
				"Schema.Educations.Entities[0].University":      "required",
				"Schema.Educations.Entities[0].StartDate":       "required",
				"Schema.Educations.Entities[0].Details[0]":      "min",
				"Schema.Educations.Entities[0].Technologies[0]": "min",
			},
		},
		{
			name: "invalid certificates entity",
			content: `template: {path: "path"}
bio: {name: "ho", title: "title"}
certificates: {entities: [{expirationDate: "2025-01-19"}]}`,
			failedValidationFields: map[string]string{
				"Schema.Certificates.Entities[0].Title": "required",
			},
		},
		{
			name: "invalid publications entity",
			content: `template: {path: "path"}
bio: {name: "ho", title: "title"}
publications: {entities: [{details: ["a"], "link": "x"}]}`,
			failedValidationFields: map[string]string{
				"Schema.Publications.Entities[0].Title":       "required",
				"Schema.Publications.Entities[0].Publisher":   "required",
				"Schema.Publications.Entities[0].PublishDate": "required",
				"Schema.Publications.Entities[0].Link":        "url",
				"Schema.Publications.Entities[0].Details[0]":  "min",
			},
		},
		{
			name: "invalid skills entity",
			content: `template: {path: "path"}
bio: {name: "ho", title: "title"}
skills: {entities: [{items: [{}]}]}`,
			failedValidationFields: map[string]string{
				"Schema.Skills.Entities[0].Category":      "required",
				"Schema.Skills.Entities[0].Items[0].Name": "required",
			},
		},
		{
			name: "invalid projects entity",
			content: `template: {path: "path"}
bio: {name: "ho", title: "title"}
projects: {entities: [{details: [""]}]}`,
			failedValidationFields: map[string]string{
				"Schema.Projects.Entities[0].Title":      "required",
				"Schema.Projects.Entities[0].Link":       "required",
				"Schema.Projects.Entities[0].Details[0]": "min",
			},
		},
		{
			name: "invalid custom sections entity",
			content: `template: {path: "path"}
bio: {name: "ho", title: "title"}
customSections: [{details: ["a"], header: ""}]`,
			failedValidationFields: map[string]string{
				"Schema.CustomSections[0].Header":     "required",
				"Schema.CustomSections[0].Details[0]": "min",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			data, err := types.NewSchema([]byte(tc.content), types.SchemaTypeYaml)

			if tc.hasInitError {
				require.Error(t, err)
				require.Nil(t, data)

				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, data)

			err = data.IsValid()

			if len(tc.failedValidationFields) > 0 {
				require.Error(t, err)

				ve := validator.ValidationErrors{}

				require.ErrorAsf(t, err, &ve, "error is not of type ValidationErrors")

				require.Len(t, ve, len(tc.failedValidationFields))

				for _, field := range ve {
					failedTag, failedFieldFound := tc.failedValidationFields[field.StructNamespace()]

					require.Truef(
						t, failedFieldFound,
						"field %s not found in failed validation fields: %v",
						field.StructNamespace(), ve,
					)

					require.Equalf(
						t, failedTag, field.Tag(),
						"field %s failed for %s instead of %s",
						field.StructNamespace(), field.Tag(), failedTag,
					)
				}

				return
			}

			require.NoError(t, err)
		})
	}
}

func TestExampleSchema(t *testing.T) {
	t.Parallel()

	data, err := types.NewSchema(sampleSchema, types.SchemaTypeYaml)

	require.NoError(t, err)

	err = data.IsValid()

	require.NoError(t, err)

	testCases := []struct {
		name     string
		key      string
		validate func(value any) bool
	}{
		{
			name: "template path",
			key:  "template.path",
			validate: func(value any) bool {
				d := value.(string)

				return d == "https://raw.githubusercontent.com/seinshah/civic/main/examples/example.template.html"
			},
		},
		{
			name: "bio contact socials",
			key:  "bio.contact.socials",
			validate: func(value any) bool {
				d := value.([]any)

				return len(d) == 2
			},
		},
		{
			name: "work experiences header",
			key:  "workExperiences.header",
			validate: func(value any) bool {
				d := value.(string)

				return d == "Work Experiences2"
			},
		},
		{
			name: "education first entity details",
			key:  "educations.entities[0].details",
			validate: func(value any) bool {
				d := value.([]any)

				return len(d) == 3
			},
		},
		{
			name: "skills default header",
			key:  "skills.header",
			validate: func(value any) bool {
				d := value.(string)

				return d == "Skills"
			},
		},
		{
			name: "second custom section details",
			key:  "customSections[1].details",
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
