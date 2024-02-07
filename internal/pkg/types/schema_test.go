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
)

const sampleSchema = `template:
  path: https://raw.githubusercontent.com/argoproj/argo-cd/master/manifests/install.yaml

bio:
  name: Hossein Shahsahebi
  title: Backend Lead Engineer
  about: |
    Solution-driven and passionate senior software engineer with more than 7 years of
    experience designing, architect, and developing high-load and well-designed server-side
    application across different industries using modern tools and technologies.
  contact:
    location: City, Country
    website: https://example.com
    email: me@example.com
    phone: +1 123 456 7890
    social:
      - https://linkedin.com/in/username/
      - https://github.com/username
  customData:
    - label: "Notice Period"
      value: "3 Months"

workExperiences:
  header: "Work Experiences2"
  entities:
    - title: Backend Lead Engineer
      company: Alphabet Inc.
      startDate: 07/2021
      endDate: present
      location: Mountain View, CA
      details:
        - Lead a team of 5 backend engineers
        - Design and develop a microservice-based system
        - Implement CI/CD pipeline
        - Conduct code reviews and pair programming
        - Mentor junior engineers
      technologies:
        - Go
        - Kubernetes
        - Docker

educations:
  header: "Educations2"
  entities:
    - degree: M	aster of Science (MSc)
      field: Computer Science
      university: University of Tehran
      startDate: 09/2014
      endDate: 09/2016
      location: Tehran, Iran
      details:
        - "Thesis: A Novel Approach for Detecting and Preventing SQL Injection Attacks"
        - "GPA: 3.8/4.0"
        - "Courses: Advanced Database, Data Mining, and Machine Learning"
      technologies:
        - Java
        - Python

certificates:
  header: "Certificates2"
  entities:
    - title: Certified Kubernetes Administrator (CKA)
      issuer: Cloud Native Computing Foundation
      issueDate: 07/2021
      expirationDate: 07/2024

publications:
  header: "Publications2"
  entities:
    - title: "A Novel Approach for Detecting and Preventing SQL Injection Attacks"
      publisher: IEEE
      publishDate: 09/2016
      link: https://ieeexplore.ieee.org/document/1234567
      details:
        - "This paper presents a novel approach for detecting and preventing SQL injection attacks using a combination of static and dynamic analysis."

skills:
  entities:
    - category: Backend
      items:
        - Name: Go
          Description: "Proficient"

projects:
  header: "Projects2"
  entities:
    - title: "Argo CD"
      link: https://argoproj.io
      details:
        - "Argo CD is a declarative, GitOps continuous delivery tool for Kubernetes."

customSections:
  - header: "Languages"
    details:
      - "some non formatted information"
  - header: "Hobbies"
    details:
      - "something about hobby 1"
      - "something about hobby 2"
      - "something about hobby 3"`

var re = regexp.MustCompile(`\[(?P<sliceIndex>\d+)]`)

func TestSchemaDataType(t *testing.T) {
	t.Parallel()

	var data types.Schema

	err := defaults.Set(&data)

	require.NoError(t, err)

	err = yaml.Unmarshal([]byte(sampleSchema), &data)

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

				return d == "https://raw.githubusercontent.com/argoproj/argo-cd/master/manifests/install.yaml"
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
