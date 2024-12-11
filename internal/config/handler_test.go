package config_test

import (
	"context"
	"os"
	"testing"

	"github.com/seinshah/cvci/internal/config"
	"github.com/seinshah/cvci/internal/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestHandler_Init(t *testing.T) {
	t.Parallel()

	h := config.NewHandler()

	wd, _ := os.Getwd()

	testCases := []struct {
		name       string
		ctx        func() (context.Context, context.CancelFunc)
		outputPath string
		hasError   bool
		err        error
	}{
		{
			name: "success",
		},
		{
			name:       "success with custom path",
			outputPath: os.TempDir() + "/test.yaml",
		},
		{
			name: "error loading sample config",
			ctx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 0)
			},
			hasError: true,
			err:      context.DeadlineExceeded,
		},
		{
			name:       "error non existing path",
			outputPath: "./non-existing-path/sample.yaml",
			hasError:   true,
		},
		{
			name:       "existing dir",
			outputPath: wd,
			hasError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			if tc.ctx != nil {
				var cancel func()

				ctx, cancel = tc.ctx()
				defer cancel()
			}

			err := h.Init(ctx, tc.outputPath)

			expectedPath := tc.outputPath
			if expectedPath == "" {
				expectedPath = "./" + types.DefaultConfigFileName
			}

			defer func() {
				_ = os.Remove(expectedPath)
			}()

			if tc.hasError {
				require.Error(t, err)

				if tc.err != nil {
					require.ErrorIs(t, err, tc.err)
				}

				return
			}

			require.NoError(t, err)

			content, err := os.ReadFile(expectedPath)

			require.NoError(t, err)
			require.NotEmpty(t, content)
			require.Greater(t, len(content), 1)
		})
	}
}
