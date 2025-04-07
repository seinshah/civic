---
slug: /configuration-structure
description: Detailed explanation of the configuration file structure.
keyword: cv, resume, curriculum vitae, civic, job, work
---

# Configuration Structure

To generate your CV, you need to provide a Civic configuration file that
contains your details. This file has to be organised in an specific way
which is explained here.

Currently, Civic only supports the YAML format for the configuration file.

| Key               | Data Type                                    | Required | Description                                                                            |
|-------------------|----------------------------------------------|----------|----------------------------------------------------------------------------------------|
| `template`        | [object(Template)](#Template)                | ✅        | template information                                                                   |
| `page`            | [object(Page)](#Page)                        | ❌        | output page setup                                                                      |
| `bio`             | [object(Bio)](#Bio)                          | ✅        | personal information                                                                   |
| `workExperiences` | [object(WorkExperiences)](#Work-Experiences) | ❌        | list of work experiences (default: empty)                                              |
| `educations`      | [object(Educations)](#Educations)            | ❌        | list of degrees (default: empty)                                                       |
| `certificates`    | [object(Certificates)](#Certificates)        | ❌        | list of certificates (default: empty)                                                  |
| `publications`    | [object(Publications)](#Publications)        | ❌        | list of publications (default: empty)                                                  |
| `skills`          | [object(Skills)](#Skills)                    | ❌        | list of skills (default: empty)                                                        |
| `projects`        | [object(Projects)](#Projects)                | ❌        | list of projects (default: empty)                                                      |
| `customSections`  | [object(CustomSections)](#Custom-Sections)   | ❌        | any additional section that doesn't fall into any of the defined ones (default: empty) |


## Template

| Key          | Data Type                        | Required | Description                                                                                         |
|--------------|----------------------------------|----------|-----------------------------------------------------------------------------------------------------|
| `path`       | string                           | ✅        | relative or absolute path to the template file on the local machine or http link to the remote file |
| `customizer` | [object(Customizer)](Customizer) | ❌        | customize template's design                                                                         |

### Customizer

| Key     | Data Type | Required | Description                        |
|---------|-----------|----------|------------------------------------|
| `style` | string    | ❌        | raw css to override template's css |

## Page

| Key      | Data Type                          | Required | Description                    |
|----------|------------------------------------|----------|--------------------------------|
| `size`   | enum(A4, B4, A, Arch-A, Letter)    | ❌        | output page size (default: A4) |
| `margin` | [object(PageMargin)](#Page-Margin) | ❌        | output page margin             |

### Page Margin

| Key      | Data Type | Required | Description                                          |
|----------|-----------|----------|------------------------------------------------------|
| `top`    | number    | ❌        | top margin in inches (default: 0, min: 0, max: 3)    |
| `right`  | number    | ❌        | right margin in inches (default: 0, min: 0, max: 3)  |
| `bottom` | number    | ❌        | bottom margin in inches (default: 0, min: 0, max: 3) |
| `left`   | number    | ❌        | left margin in inches (default: 0, min: 0, max: 3)   |

## Bio

| Key          | Data Type                          | Required | Description                                                                           |
|--------------|------------------------------------|----------|---------------------------------------------------------------------------------------|
| `name`       | string                             | ✅        | full name (max length: 100)                                                           |
| `title`      | string                             | ✅        | current occupation title                                                              |
| `about`      | string                             | ❌        | a paragraph about your experience and achievements that usually appears on top of CVs |
| `contact`    | [object(Contact)](#Contact)        | ❌        | your contact details (default: empty)                                                 |
| `customData` | [object(CustomData)](#Custom-Data) | ❌        | custom bio sections that don't fall in any other category (default: empty)            |

### Contact

| Key        | Data Type     | Required | Description                       |
|------------|---------------|----------|-----------------------------------|
| `location` | string        | ❌        | your residence                    |
| `website`  | string        | ❌        | link to your personal website     |
| `email`    | string        | ✅        | email address                     |
| `phone`    | string        | ❌        | phone number                      |
| `socials`  | array(string) | ❌        | list of social media profile URLs |

### Custom Data

| Key     | Data Type | Required | Description                                   |
|---------|-----------|----------|-----------------------------------------------|
| `label` | string    | ❌        | label of the custom bio info (default: empty) |
| `value` | string    | ✅        | phrase to be shown as custom bio detail       |

## Work Experiences

| Key        | Data Type                                                      | Required | Description                                 |
|------------|----------------------------------------------------------------|----------|---------------------------------------------|
| `header`   | string                                                         | ❌        | section title (default: "Work Experiences") |
| `entities` | [array(object(WorkExperienceEntity))](#Work-Experience-Entity) | ✅        | list of experiences (min items: 1)          |

### Work Experience Entity

| Key            | Data Type     | Required | Description                             |
|----------------|---------------|----------|-----------------------------------------|
| `title`        | string        | ✅        | job title                               |
| `company`      | string        | ✅        | company name                            |
| `location`     | string        | ❌        | job location                            |
| `startDate`    | string        | ✅        | start date of the job                   |
| `endDate`      | string        | ❌        | end date of the job (default: present)  |
| `details`      | array(string) | ❌        | itemized description of your activities |
| `technologies` | array(string) | ❌        | list of techs you worked with           |

## Educations

| Key        | Data Type                                           | Required | Description                           |
|------------|-----------------------------------------------------|----------|---------------------------------------|
| `header`   | string                                              | ❌        | section title (default: "Educations") |
| `entities` | [array(object(EducationEntity))](#Education-Entity) | ✅        | list of degrees (min items: 1)        |

### Education Entity

| Key            | Data Type     | Required | Description                                   |
|----------------|---------------|----------|-----------------------------------------------|
| `degree`       | string        | ✅        | degree title                                  |
| `field`        | string        | ✅        | field of study                                |
| `university`   | string        | ✅        | university name                               |
| `location`     | string        | ❌        | university location                           |
| `startDate`    | string        | ✅        | start date of study                           |
| `endDate`      | string        | ❌        | end date of study (default: present)          |
| `details`      | array(string) | ❌        | itemized list of interesting details to share |
| `technologies` | array(string) | ❌        | list of technologies you worked with          |

## Certificates

| Key        | Data Type                                               | Required | Description                             |
|------------|---------------------------------------------------------|----------|-----------------------------------------|
| `header`   | string                                                  | ❌        | section title (default: "Certificates") |
| `entities` | [array(object(CertificateEntity))](#Certificate-Entity) | ✅        | list of certificates (min items: 1)     |

### Certificate Entity

| Key              | Data Type | Required | Description                                           |
|------------------|-----------|----------|-------------------------------------------------------|
| `title`          | string    | ✅        | certificate title                                     |
| `issuer`         | string    | ✅        | certificate issuer                                    |
| `issueDate`      | string    | ✅        | issue date of the certificate                         |
| `expirationDate` | string    | ❌        | expiration date of the certificate (default: present) |

## Publications

| Key        | Data Type                                               | Required | Description                             |
|------------|---------------------------------------------------------|----------|-----------------------------------------|
| `header`   | string                                                  | ❌        | section title (default: "Publications") |
| `entities` | [array(object(PublicationEntity))](#Publication-Entity) | ✅        | list of publications (min items: 1)     |

### Publication Entity

| Key           | Data Type     | Required | Description                            |
|---------------|---------------|----------|----------------------------------------|
| `title`       | string        | ✅        | publication title                      |
| `publisher`   | string        | ✅        | publisher name                         |
| `publishDate` | string        | ✅        | publication date                       |
| `link`        | string        | ✅        | http link to the published item        |
| `details`     | array(string) | ❌        | itemized details about the publication |

## Skills

| Key        | Data Type                                   | Required | Description                       |
|------------|---------------------------------------------|----------|-----------------------------------|
| `header`   | string                                      | ❌        | section title (default: "Skills") |
| `entities` | [array(object(SkillEntity))](#Skill-Entity) | ✅        | list of skills (min items: 1)     |

### Skill Entity

| Key        | Data Type                              | Required | Description                                   |
|------------|----------------------------------------|----------|-----------------------------------------------|
| `category` | string                                 | ✅        | category name to group list of skills         |
| `items`    | [array(object(SkillItem))[#Skill-Item] | ✅        | list of skills in the category (min items: 1) |

### Skill Item

| Key           | Data Type | Required | Description                                                                                           |
|---------------|-----------|----------|-------------------------------------------------------------------------------------------------------|
| `name`        | string    | ✅        | skill name                                                                                            |
| `description` | string    | ❌        | description of the skill (this might be used by template files differently for more appealing design) |

## Projects

| Key        | Data Type                                       | Required | Description                         |
|------------|-------------------------------------------------|----------|-------------------------------------|
| `header`   | string                                          | ❌        | section title (default: "Projects") |
| `entities` | [array(object(ProjectEntity))](#Project-Entity) | ✅        | list of projects (min items: 1)     |

### Project Entity

| Key       | Data Type     | Required | Description                                 |
|-----------|---------------|----------|---------------------------------------------|
| `title`   | string        | ✅        | project title                               |
| `link`    | string        | ✅        | http link to the project                    |
| `details` | array(string) | ❌        | itemized details to share about the project |

## Custom Sections

| Key       | Data Type     | Required | Description                                    |
|-----------|---------------|----------|------------------------------------------------|
| `header`  | string        | ✅        | section title                                  |
| `details` | array(string) | ✅        | itemized list of details in the custom section |
