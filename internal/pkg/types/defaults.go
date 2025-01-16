package types

import "os"

const (
	DefaultOutputFileName = "cvci.pdf"
	DefaultSchemaFileName = ".cvci.yaml"
	DefaultPageSize       = PageSizeA4
	DefaultFilePermission = 0o600
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
