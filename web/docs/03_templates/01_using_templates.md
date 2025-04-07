---
sidebar_position: 1
---

# Using Templates

Templates in Civic allow you to customize the appearance and structure of your CV/Resume. This guide will help you understand how to choose and use templates effectively.

## Understanding Templates

A template in Civic is an HTML file that defines the structure and styling of your CV/Resume. Each template:
- Must be compatible with your Civic version
- Can include custom CSS styling
- Supports various sections defined in your CV schema
- Can be customized further through the configuration

## Choosing a Template

When selecting a template, consider the following factors:

1. **Version Compatibility**: Templates are version-specific. Make sure the template's version matches your Civic installation's major version.

2. **Direction Support**: Templates can support different text directions (LTR/RTL) for international usage.

3. **Style Preferences**: Each template comes with its own default styling, which you can preview before using.

## Template Configuration

To use a template, you need to specify it in your CV configuration file:

```yaml
template:
  path: "./path/to/template.html"  # Path to your template file
  customizer:
    style: |
      /* Optional: Add custom CSS styles */
      body {
        font-family: 'Arial', sans-serif;
      }
```

### Customization Options

You can customize templates in two ways:

1. **Direct Template Selection**: Choose a different template file
2. **CSS Customization**: Add custom CSS through the `customizer.style` property

## Template Validation

Civic automatically validates templates to ensure:
- They are compatible with your Civic version
- They don't contain any forbidden HTML tags (for security)
- All required meta information is present

## Best Practices

1. **Version Check**: Always check the template's version compatibility before using it
2. **Preview First**: Test the template with a sample CV before using it for your final document
3. **Backup**: Keep a backup of your working template configuration
4. **Custom Styles**: Start with minimal custom CSS and add styles incrementally as needed

## Troubleshooting

Common template-related issues and solutions:

1. **Version Mismatch**: If you see a version mismatch error, check the template's meta tags and ensure they match your Civic version
2. **Missing Sections**: Ensure your template includes placeholders for all sections in your CV schema
3. **Styling Issues**: Use the browser's developer tools to inspect and debug styling problems
