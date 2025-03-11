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
	Style string `json:"style,omitempty" yaml:"style"`
}

type SchemaTemplate struct {
	// Path is the path to the template file. It can be a local path to the template file
	// on the host, or an HTTP link to where the template is located.
	// If you provide a remote path, the link should refer to the raw HTML file.
	Path string `json:"path" validate:"required" yaml:"path"`

	// Customizer is a way for you to customize the template in use.
	Customizer Customizer `json:"customizer,omitempty" yaml:"customizer"`
}

type SchemaPage struct {
	// Size is the size of the page for the PDF.
	// Valid values are: A4, B4, A, Arch-A, Letter.
	// If an invalid value is provided, it will default to A4.
	Size PageSize `default:"A4" json:"size,omitempty" yaml:"size"`

	// Margin is the margin of the page for the PDF.
	// Absence of margin for each side leads to 0.
	// IMPORTANT: Dimensions are in inch.
	Margin PageMargin `json:"margin,omitempty" yaml:"margin"`
}

type SchemaBioContact struct {
	// Location is the current location of the person.
	Location string `json:"location,omitempty" yaml:"location"`

	// Website contains the link to the portfolio of the person.
	Website string `json:"website,omitempty" validate:"omitempty,url" yaml:"website"`

	// Email is the email address of the person.
	Email string `json:"email" validate:"required,email" yaml:"email"`

	// Phone is the phone number of the person.
	Phone string `json:"phone,omitempty" yaml:"phone"`

	// Social contains the links to the social media profiles of the person.
	Socials []string `json:"socials" validate:"dive,url" yaml:"socials"`
}

type SchemaBioCustomData struct {
	// Label or title of the custom data.
	Label string `json:"label,omitempty" yaml:"label"`

	// Value of the custom data.
	Value string `json:"value" validate:"required" yaml:"value"`
}

type SchemaBio struct {
	// Name is the full name of the person.
	Name string `json:"name" validate:"required,min=2,max=100" yaml:"name"`

	// Title is the career title of the person.
	Title string `json:"title" validate:"required,min=2" yaml:"title"`

	// About is a short description about the person.
	About string `json:"about,omitempty" yaml:"about"`

	Contact *SchemaBioContact `json:"contact,omitempty" validate:"omitempty" yaml:"contact"`

	// CustomData is the list of any additional key and values that is going to be
	// part of the personal information section of the resume or cv.
	CustomData []SchemaBioCustomData `json:"customData,omitempty" validate:"omitempty,dive" yaml:"customData"`
}

type SchemaWorkExperienceEntity struct {
	// Title is the title of the job.
	Title string `json:"title" validate:"required" yaml:"title"`

	// Company is the name of the company.
	Company string `json:"company" validate:"required" yaml:"company"`

	// Location is the location of the job.
	Location string `json:"location,omitempty" yaml:"location"`

	// StartDate is the start date of the job. There is no validation for the date format.
	StartDate string `json:"startDate" validate:"required" yaml:"startDate"`

	// EndDate is the end date of the job. There is no validation for the date format.
	EndDate string `default:"present" json:"endDate,omitempty" yaml:"endDate"`

	// Details is the list of details about the job. There is no validation.
	// It can include the list of achievements, responsibilities, and any other details.
	Details []string `json:"details,omitempty" validate:"dive,min=2" yaml:"details"`

	// Technologies are the list of tools and technologies that you were exposed to during the job.
	Technologies []string `json:"technologies,omitempty" validate:"dive,min=1" yaml:"technologies"`
}

type SchemaWorkExperiences struct {
	// Header is the printed header/title of this section.
	Header string `default:"Work Experiences" json:"header,omitempty" yaml:"header"`

	Entities []SchemaWorkExperienceEntity `json:"entities" validate:"required,min=1,dive" yaml:"entities"`
}

type SchemaEducationsEntity struct {
	// Degree is the degree that you have achieved.
	Degree string `json:"degree" validate:"required" yaml:"degree"`

	// Field is the field of study.
	Field string `json:"field" validate:"required" yaml:"field"`

	// University is the name of the university or place of study.
	University string `json:"university" validate:"required" yaml:"university"`

	// Location is the location of the university or place of study.
	Location string `json:"location,omitempty" yaml:"location"`

	// StartDate is the start date of the study. There is no validation for the date format.
	StartDate string `json:"startDate" validate:"required" yaml:"startDate"`

	// EndDate is the end date of the study. There is no validation for the date format.
	EndDate string `default:"present" json:"endDate,omitempty" yaml:"endDate"`

	// Details is the list of details about the study. There is no validation.
	// It can include the list of achievements, responsibilities, and any other details.
	Details []string `json:"details,omitempty" validate:"dive,min=2" yaml:"details"`

	// Technologies are the list of tools and technologies that you were exposed to during the study.
	Technologies []string `json:"technologies,omitempty" validate:"dive,min=1" yaml:"technologies"`
}

type SchemaEducations struct {
	// Header is the printed header/title of this section.
	Header string `default:"Educations" json:"header,omitempty" yaml:"header"`

	Entities []SchemaEducationsEntity `json:"entities" validate:"required,min=1,dive" yaml:"entities"`
}

type SchemaCertificatesEntity struct {
	// Title is the title of the certificate.
	Title string `json:"title" validate:"required" yaml:"title"`

	// Issuer is the name of the issuer of the certificate.
	Issuer string `json:"issuer" validate:"required" yaml:"issuer"`

	// IssueDate is the date when the certificate was issued. There is no validation for the date format.
	IssueDate string `json:"issueDate" validate:"required" yaml:"issueDate"`

	// ExpiryDate is the date when the certificate will expire. There is no validation for the date format.
	ExpirationDate string `json:"expirationDate" yaml:"expirationDate"`
}

type SchemaCertificates struct {
	// Header is the printed header/title of this section.
	Header string `default:"Certificates" json:"header,omitempty" yaml:"header"`

	Entities []SchemaCertificatesEntity `json:"entities" validate:"required,min=1,dive" yaml:"entities"`
}

type SchemaPublicationsEntity struct {
	// Title is the title of the publication.
	Title string `json:"title" validate:"required" yaml:"title"`

	// Publisher is the name of the publisher of the publication.
	Publisher string `json:"publisher" validate:"required" yaml:"publisher"`

	// PublishDate is the date when the publication was published. There is no validation for the date format.
	PublishDate string `json:"publishDate" validate:"required" yaml:"publishDate"`

	// Link is the link to the publication.
	Link string `json:"link" validate:"required,url" yaml:"link"`

	// Details is the list of details about the publication. There is no validation.
	Details []string `json:"details,omitempty" validate:"dive,min=2" yaml:"details"`
}

type SchemaPublications struct {
	// Header is the printed header/title of this section.
	Header string `default:"Publications" json:"header,omitempty" yaml:"header"`

	Entities []SchemaPublicationsEntity `json:"entities" validate:"required,min=1,dive" yaml:"entities"`
}

type SchemaSkillsEntityItem struct {
	// Name is the name of the skill.
	Name string `json:"name" validate:"required" yaml:"name"`

	// Description is any arbitrary detail of this skill. This might be interpreted
	// in certain way by each template.
	Description string `json:"description,omitempty" yaml:"description"`
}

type SchemaSkillsEntity struct {
	// Category is the category of the skill.
	Category string `json:"category" validate:"required" yaml:"category"`

	// Items contain all the tools and technologies in this category that you have experience with.
	Items []SchemaSkillsEntityItem `json:"items" validate:"required,min=1,dive" yaml:"items"`
}

type SchemaSkills struct {
	// Header is the printed header/title of this section.
	Header string `default:"Skills" json:"header,omitempty" yaml:"header"`

	Entities []SchemaSkillsEntity `json:"entities" validate:"required,min=1,dive" yaml:"entities"`
}

type SchemaProjectsEntity struct {
	// Title is the title of the project.
	Title string `json:"title" validate:"required" yaml:"title"`

	// Link is the link to the project.
	Link string `json:"link" validate:"required,url" yaml:"link"`

	// Details is the list of details about the project.
	Details []string `json:"details,omitempty" validate:"dive,min=1" yaml:"details"`
}

type SchemaProjects struct {
	// Header is the printed header/title of this section.
	Header string `default:"Projects" json:"header,omitempty" yaml:"header"`

	Entities []SchemaProjectsEntity `json:"entities" validate:"required,min=1,dive" yaml:"entities"`
}

type SchemaCustomSection struct {
	// Header is the title of the custom section.
	Header string `json:"header" validate:"required,min=1" yaml:"header"`

	// A list of arbitrary details to be shown under this section.
	Details []string `json:"details" validate:"required,min=1,dive,min=2" yaml:"details"`
}

// Schema is the architecture of the configuration file that will be provided
// by the user to be used for generating the final resume or cv.
type Schema struct {
	// Template contains all the information related to the template file
	// that will be used to create the resume or cv.
	Template SchemaTemplate `json:"template,omitempty" validate:"required" yaml:"template"`

	Page SchemaPage `json:"page,omitempty" yaml:"page"`

	// Bio contains all the personal information of the person.
	Bio SchemaBio `json:"bio" validate:"required" yaml:"bio"`

	// WorkExperiences contains all the work experiences of the person.
	WorkExperiences *SchemaWorkExperiences `json:"workExperiences,omitempty" validate:"omitempty" yaml:"workExperiences"`

	// Educations contains all the educations of the person.
	Educations *SchemaEducations `json:"educations,omitempty" validate:"omitempty" yaml:"educations"`

	// Certificates contains all the certificates of the person.
	Certificates *SchemaCertificates `json:"certificates,omitempty" validate:"omitempty" yaml:"certificates"`

	// Publications contains all the publications of the person.
	Publications *SchemaPublications `json:"publications,omitempty" validate:"omitempty" yaml:"publications"`

	// Skills contains all the skills of the person separated by category.
	Skills *SchemaSkills `json:"skills,omitempty" validate:"omitempty" yaml:"skills"`

	// Projects contains all the projects of the person.
	Projects *SchemaProjects `json:"projects,omitempty" validate:"omitempty" yaml:"projects"`

	// CustomSections contains all the custom sections that you want to add to the resume or cv.
	CustomSections []SchemaCustomSection `json:"customSections,omitempty" validate:"omitempty,dive" yaml:"customSections"`
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
