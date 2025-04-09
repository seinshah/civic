package output_test

import (
	"os"
	"testing"

	"github.com/seinshah/civic/internal/pkg/output"
	"github.com/seinshah/civic/internal/pkg/output/html"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	t.Parallel()

	engine := html.NewEngine()
	content := []byte("<p>test content</p>")

	err := output.Render(t.Context(), content, engine, "/tmp/test.html")

	require.NoError(t, err)

	require.FileExists(t, "/tmp/test.html")

	defer func() {
		_ = os.Remove("/tmp/test.html")
	}()

	actualContent, err := os.ReadFile("/tmp/test.html")

	require.NoError(t, err)
	require.Equal(t, content, actualContent)
}
