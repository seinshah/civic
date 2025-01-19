// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package types

import (
	"fmt"
	"strings"
)

const (
	// SchemaTypeYaml is a SchemaType of type yaml.
	SchemaTypeYaml SchemaType = "yaml"
	// SchemaTypeYml is a SchemaType of type yml.
	SchemaTypeYml SchemaType = "yml"
)

var ErrInvalidSchemaType = fmt.Errorf("not a valid SchemaType, try [%s]", strings.Join(_SchemaTypeNames, ", "))

var _SchemaTypeNames = []string{
	string(SchemaTypeYaml),
	string(SchemaTypeYml),
}

// SchemaTypeNames returns a list of possible string values of SchemaType.
func SchemaTypeNames() []string {
	tmp := make([]string, len(_SchemaTypeNames))
	copy(tmp, _SchemaTypeNames)
	return tmp
}

// String implements the Stringer interface.
func (x SchemaType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x SchemaType) IsValid() bool {
	_, err := ParseSchemaType(string(x))
	return err == nil
}

var _SchemaTypeValue = map[string]SchemaType{
	"yaml": SchemaTypeYaml,
	"yml":  SchemaTypeYml,
}

// ParseSchemaType attempts to convert a string to a SchemaType.
func ParseSchemaType(name string) (SchemaType, error) {
	if x, ok := _SchemaTypeValue[name]; ok {
		return x, nil
	}
	return SchemaType(""), fmt.Errorf("%s is %w", name, ErrInvalidSchemaType)
}
