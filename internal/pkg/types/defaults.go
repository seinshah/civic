package types

import "os"

const (
	DefaultAppName            = "civic"
	DefaultAppOwner           = "seinshah"
	DefaultOutputFileName     = DefaultAppName + ".pdf"
	DefaultSchemaFileName     = "." + DefaultAppName + ".yaml"
	DefaultSchemaJSONFileName = DefaultAppName + "-jsonschema.json"
	DefaultPageSize           = PageSizeA4
	DefaultFilePermission     = 0o600
)

func CurrentWDPath(filename string) string {
	workingDir, err := os.Getwd()
	if err != nil {
		workingDir = "."
	}

	if filename == "" {
		return workingDir
	}

	return workingDir + string(os.PathSeparator) + filename
}
