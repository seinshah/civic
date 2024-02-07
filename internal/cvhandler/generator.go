package cvhandler

import (
	"github.com/seinshah/cvci/internal/pkg/types"
)

type Generator struct {
	resumeConfigPath string
}

var _ types.Generator = (*Generator)(nil)

func NewGenerator(configPath string) *Generator {
	return &Generator{
		resumeConfigPath: configPath,
	}
}

func (*Generator) Generate() error {
	return nil
}
