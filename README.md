# Civic - CV as a Code

Maintaining a CV should not be a hassle. Civic is here to make sure of that. 
It splits your CV into its content and its design to help you change each 
independently.

## Installation

### From Source

Use Go to install the package directly from the source code. This will 
install the latest version of the tool. Therefore, make sure your template 
matches this version.

```bash
go install github.com/seinshah/civic@latest
```

## Usage

```
A tool to help maintaining and extending CVs or resumes 
easily by separating the template from the content.

Usage:
  civic [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    Generate the resume or cv
  help        Help about any command
  schema      Helper commands to work with CV schema file.

Flags:
  -h, --help       help for civic
      --no-color   disable the color in the logs output
      --verbose    show more information during the process
  -v, --version    version for civic

Use "civic [command] --help" for more information about a command.
```

## Schema File

The file where you define your CV content is called the schema file. It is 
following a specific json schema. Schema file can be crafted in the 
following formats:

- YAML

## Templates
Templates are design you can use to render your CV into different formats. 
The following formats are supported:

- HTML
- PDF

Templates are versioned to match the app version they are compatible with. 
However only the major versions are considered. So a template might not 
support all the features of the recent app version, but it will still work 
(`v0` might be an exception).

If you like a design, but prefer a tiny change, you don't need to craft your 
own template. The schema file allows you to customize your chosen templates 
using simple CSS directives.