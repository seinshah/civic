package types

import (
	"errors"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

var ErrEmptySchemaPath = errors.New("schema file path is empty")

//go:generate go-enum --names

// SchemaType is the type of schema being used to generate the CV.
// ENUM(yaml, yml).
type SchemaType string

type Customizer struct {
	// Style is a block of css code that will be added in a style tag
	// at the end of the HEAD section of the template.
	Style string `yaml:"style"`
}

// Schema is the architecture of the configuration file that will be provided
// by the user to be used for generating the final resume or cv.
type Schema struct {
	// Template contains all the information related to the template file
	// that will be used to create the resume or cv.
	Template struct {
		// Path is the path to the template file. It can be a local path to the template file
		// on the host, or an HTTP link to where the template is located.
		// If you provide a remote path, the link should refer to the raw HTML file.
		Path string `validate:"required" yaml:"path"`

		// Customizer is a way for you to customize the template in use.
		Customizer Customizer `yaml:"customizer"`
	} `validate:"required" yaml:"template"`

	Page struct {
		// Size is the size of the page for the PDF.
		// Valid values are: A4, B4, A, Arch-A, Letter.
		// If an invalid value is provided, it will default to A4.
		Size PageSize `default:"A4" yaml:"size"`

		// Margin is the margin of the page for the PDF.
		// Absence of margin for each side leads to 0.
		// IMPORTANT: Dimensions are in inch.
		Margin PageMargin `yaml:"margin"`
	} `yaml:"page"`

	// Bio contains all the personal information of the person.
	Bio struct {
		// Name is the full name of the person.
		Name string `validate:"required,min=2,max=100" yaml:"name"`

		// Title is the career title of the person.
		Title string `validate:"required,min=2" yaml:"title"`

		// About is a short description about the person.
		About string `yaml:"about"`

		Contact *struct {
			// Location is the current location of the person.
			Location string `yaml:"location"`

			// Website contains the link to the portfolio of the person.
			Website string `validate:"omitempty,url" yaml:"website"`

			// Email is the email address of the person.
			Email string `validate:"required,email" yaml:"email"`

			// Phone is the phone number of the person.
			Phone string `yaml:"phone"`

			// Social contains the links to the social media profiles of the person.
			Socials []string `validate:"dive,url" yaml:"social"`
		} `validate:"omitempty" yaml:"contact"`

		// CustomData is the list of any additional key and values that is going to be
		// part of the personal information section of the resume or cv.
		CustomData []struct {
			// Label or title of the custom data.
			Label string `yaml:"label"`

			// Value of the custom data.
			Value string `validate:"required" yaml:"value"`
		} `validate:"omitempty,dive" yaml:"customData"`
	} `validate:"required" yaml:"bio"`

	// WorkExperiences contains all the work experiences of the person.
	WorkExperiences *struct {
		// Header is the printed header/title of this section.
		Header string `default:"Work Experiences" yaml:"header"`

		Entities []struct {
			// Title is the title of the job.
			Title string `validate:"required" yaml:"title"`

			// Company is the name of the company.
			Company string `validate:"required" yaml:"company"`

			// Location is the location of the job.
			Location string `yaml:"location"`

			// StartDate is the start date of the job. There is no validation for the date format.
			StartDate string `validate:"required" yaml:"startDate"`

			// EndDate is the end date of the job. There is no validation for the date format.
			EndDate string `default:"present" yaml:"endDate"`

			// Details is the list of details about the job. There is no validation.
			// It can include the list of achievements, responsibilities, and any other details.
			Details []string `validate:"dive,min=2" yaml:"details"`

			// Technologies are the list of tools and technologies that you were exposed to during the job.
			Technologies []string `validate:"dive,min=1" yaml:"technologies"`
		} `validate:"required,min=1,dive" yaml:"entities"`
	} `validate:"omitempty" yaml:"workExperiences"`

	// Educations contains all the educations of the person.
	Educations *struct {
		// Header is the printed header/title of this section.
		Header string `default:"Educations" yaml:"header"`

		Entities []struct {
			// Degree is the degree that you have achieved.
			Degree string `validate:"required" yaml:"degree"`

			// Field is the field of study.
			Field string `validate:"required" yaml:"field"`

			// University is the name of the university or place of study.
			University string `validate:"required" yaml:"university"`

			// Location is the location of the university or place of study.
			Location string `yaml:"location"`

			// StartDate is the start date of the study. There is no validation for the date format.
			StartDate string `validate:"required" yaml:"startDate"`

			// EndDate is the end date of the study. There is no validation for the date format.
			EndDate string `default:"present" yaml:"endDate"`

			// Details is the list of details about the study. There is no validation.
			// It can include the list of achievements, responsibilities, and any other details.
			Details []string `validate:"dive,min=2" yaml:"details"`

			// Technologies are the list of tools and technologies that you were exposed to during the study.
			Technologies []string `validate:"dive,min=1" yaml:"technologies"`
		} `validate:"required,min=1,dive" yaml:"entities"`
	} `validate:"omitempty" yaml:"educations"`

	// Certificates contains all the certificates of the person.
	Certificates *struct {
		// Header is the printed header/title of this section.
		Header string `default:"Certificates" yaml:"header"`

		Entities []struct {
			// Title is the title of the certificate.
			Title string `validate:"required" yaml:"title"`

			// Issuer is the name of the issuer of the certificate.
			Issuer string `validate:"required" yaml:"issuer"`

			// IssueDate is the date when the certificate was issued. There is no validation for the date format.
			IssueDate string `validate:"required" yaml:"issueDate"`

			// ExpiryDate is the date when the certificate will expire. There is no validation for the date format.
			ExpirationDate string `yaml:"expirationDate"`
		} `validate:"required,min=1,dive" yaml:"entities"`
	} `validate:"omitempty" yaml:"certificates"`

	// Publications contains all the publications of the person.
	Publications *struct {
		// Header is the printed header/title of this section.
		Header string `default:"Publications" yaml:"header"`

		Entities []struct {
			// Title is the title of the publication.
			Title string `validate:"required" yaml:"title"`

			// Publisher is the name of the publisher of the publication.
			Publisher string `validate:"required" yaml:"publisher"`

			// PublishDate is the date when the publication was published. There is no validation for the date format.
			PublishDate string `validate:"required" yaml:"publishDate"`

			// Link is the link to the publication.
			Link string `validate:"required,url" yaml:"link"`

			// Details is the list of details about the publication. There is no validation.
			Details []string `validate:"dive,min=2" yaml:"details"`
		} `validate:"required,min=1,dive" yaml:"entities"`
	} `validate:"omitempty" yaml:"publications"`

	// Skills contains all the skills of the person separated by category.
	Skills *struct {
		// Header is the printed header/title of this section.
		Header string `default:"Skills" yaml:"header"`

		Entities []struct {
			// Category is the category of the skill.
			Category string `validate:"required" yaml:"category"`

			// Items contain all the tools and technologies in this category that you have experience with.
			Items []struct {
				// Name is the name of the skill.
				Name string `validate:"required" yaml:"name"`

				// Description is any arbitrary detail of this skill. This might be interpreted
				// in certain way by each template.
				Description string `yaml:"description"`
			} `validate:"required,min=1,dive" yaml:"items"`
		} `validate:"required,min=1,dive" yaml:"entities"`
	} `validate:"omitempty" yaml:"skills"`

	// Projects contains all the projects of the person.
	Projects *struct {
		// Header is the printed header/title of this section.
		Header string `default:"Projects" yaml:"header"`

		Entities []struct {
			// Title is the title of the project.
			Title string `validate:"required" yaml:"title"`

			// Link is the link to the project.
			Link string `validate:"required,url" yaml:"link"`

			// Details is the list of details about the project.
			Details []string `validate:"dive,min=1" yaml:"details"`
		} `validate:"required,min=1,dive" yaml:"entities"`
	} `validate:"omitempty" yaml:"projects"`

	// CustomSections contains all the custom sections that you want to add to the resume or cv.
	CustomSections []struct {
		// Header is the title of the custom section.
		Header string `validate:"required,min=1" yaml:"header"`

		// A list of arbitrary details to be shown under this section.
		Details []string `validate:"required,min=1,dive,min=2" yaml:"details"`
	} `validate:"omitempty,dive" yaml:"customSections"`
}

func NewSchema(content []byte, contentType SchemaType) (*Schema, error) {
	data := Schema{}

	switch contentType {
	case SchemaTypeYaml, SchemaTypeYml:
		if err := yaml.Unmarshal(content, &data); err != nil {
			return nil, err
		}

	default:
		return nil, ErrInvalidSchemaType
	}

	// defaults should be set after loading the provided schema as defaults won't be set
	// for nil structs.
	if err := defaults.Set(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *Schema) IsValid() error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(s)
}
