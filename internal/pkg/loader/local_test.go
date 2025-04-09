package loader_test

import (
	"os"
	"testing"

	"github.com/seinshah/civic/internal/pkg/loader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalLoader_Load(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		expectError
		name            string
		path            func(t *testing.T) (string, func(t *testing.T))
		expectedContent []byte
	}{
		{
			name:            "success",
			path:            localFileSetUp("success-local-test"),
			expectedContent: []byte("success-local-test"),
		},
		{
			name:            "empty-file",
			path:            localFileSetUp(""),
			expectedContent: make([]byte, 0),
		},
		{
			name: "file-not-found",
			path: func(t *testing.T) (string, func(t *testing.T)) {
				t.Helper()

				return "not-found", func(_ *testing.T) {}
			},
			expectError: expectError{
				hasError: true,
				err:      loader.ErrInvalidLocalPath,
			},
		},
		{
			name: "directory",
			path: func(t *testing.T) (string, func(t *testing.T)) {
				t.Helper()

				return os.TempDir(), func(_ *testing.T) {}
			},
			expectError: expectError{
				hasError: true,
				err:      loader.ErrInvalidLocalPath,
			},
		},
		{
			name: "non-permitted-file",
			path: func(t *testing.T) (string, func(t *testing.T)) {
				t.Helper()

				f, err := os.CreateTemp(os.TempDir(), "app-test-*.txt")

				require.NoError(t, err)

				_, err = f.WriteString("something")

				require.NoError(t, err)

				err = f.Chmod(0)

				require.NoError(t, err)

				return f.Name(), func(t *testing.T) {
					t.Helper()
					errD := f.Close()

					require.NoError(t, errD)
				}
			},
			expectError: expectError{
				hasError: true,
				err:      loader.ErrInvalidLocalPath,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			path, cleanup := tc.path(t)
			defer cleanup(t)

			l := loader.NewLocalLoader(path)
			content, err := l.Load(t.Context())

			if tc.hasError {
				require.Error(t, err)

				if tc.err != nil {
					require.ErrorIs(t, err, tc.err)
				}

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedContent, content)
		})
	}
}
