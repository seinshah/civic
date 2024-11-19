package html_test

import (
	"context"
	"testing"

	"github.com/seinshah/cvci/internal/pkg/output/html"
	"github.com/stretchr/testify/require"
)

func TestEngine_Generate(t *testing.T) {
	t.Parallel()

	content := []byte("<p>test content</p>")

	engine := html.NewEngine()

	output, err := engine.Generate(context.Background(), content)

	require.NoError(t, err)
	require.Equal(t, content, output)
}
