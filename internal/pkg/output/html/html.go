package html

import (
	"context"

	"github.com/seinshah/cvci/internal/pkg/types"
)

type Engine struct{}

var _ types.OutputGenerator = &Engine{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e Engine) Generate(_ context.Context, content []byte) ([]byte, error) {
	return content, nil
}
