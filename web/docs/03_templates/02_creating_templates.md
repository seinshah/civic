---
sidebar_position: 2
---

# Creating Templates

This guide is for developers who want to create custom templates for Civic. It covers the technical aspects of template creation, validation, and best practices.

## Template Structure

A Civic template is an HTML file that must follow these requirements:

1. **Meta Information**: Include required meta tags:
```html
<meta name="app-version" content="1.0.0" />
<meta name="template-direction" content="ltr" />
```

2. **Security Constraints**: Certain HTML tags are forbidden for security:
- `<script>` tags are not allowed
- `<iframe>` tags are not allowed
- `<link>` tags are only allowed for stylesheets

## Template Variables

Templates use Go's HTML template syntax. You have access to the entire CV schema in your templates:

```html
<div class="profile">
  <h1>{{ .Profile.FullName }}</h1>
  <p>{{ .Profile.Title }}</p>
</div>

<div class="experience">
  {{ range .Experience }}
    <div class="job">
      <h3>{{ .Position }} at {{ .Company }}</h3>
      <p>{{ .Description }}</p>
    </div>
  {{ end }}
</div>
```

## Template Validation

Civic performs several validations on templates:

1. **Version Compatibility**:
```go
// The template version must match the app's major version
if appV.Major() != templateV.Major() {
    return fmt.Errorf("expected v%d template version: %w", appV.Major(), ErrMismatchAppVersion)
}
```

2. **Security Checks**:
```go
// Forbidden tags are checked
forbiddenTags = []string{"script", "iframe", "link"}

// Some tags have exceptions
forbiddenException = map[string]map[string][]string{
    "link": {
        "rel": []string{"stylesheet"},
    },
}
```

## Customization Support

Templates should be designed to support customization through:

1. **CSS Variables**: Define key styling properties as CSS variables:
```css
:root {
  --cv-primary-color: #333;
  --cv-font-family: 'Arial', sans-serif;
  --cv-spacing-unit: 1rem;
}
```

2. **Custom Style Injection**: Templates automatically receive custom styles from the configuration:
```html
<head>
  <style type="text/css">
    /* Base styles */
  </style>
  <!-- Custom styles will be injected here -->
</head>
```

## Development Best Practices

1. **Version Management**:
- Use semantic versioning for your templates
- Document version compatibility clearly
- Test with multiple Civic versions

2. **HTML Structure**:
- Use semantic HTML elements
- Provide clear class names for styling
- Include print-friendly styles
- Support both LTR and RTL layouts

3. **Testing**:
- Test with various CV schemas
- Validate HTML structure
- Check print layout
- Test with different browsers

4. **Documentation**:
- Document template requirements
- List supported CV schema sections
- Provide customization examples
- Include sample configurations

## Error Handling

Templates should gracefully handle:

1. **Missing Data**:
```html
{{ if .Profile.Summary }}
  <div class="summary">{{ .Profile.Summary }}</div>
{{ end }}
```

2. **Empty Lists**:
```html
{{ if .Skills }}
  <div class="skills">
    {{ range .Skills }}
      <span class="skill">{{ . }}</span>
    {{ end }}
  </div>
{{ end }}
```

## Template Distribution

When distributing your template:

1. **Package Contents**:
- Template HTML file
- README with usage instructions
- Sample configuration
- Preview images

2. **Version Information**:
- Clear version number
- Compatibility requirements
- Changelog

3. **License**:
- Include license information
- Document any attribution requirements
- Specify usage restrictions
