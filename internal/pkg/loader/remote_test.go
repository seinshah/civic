package loader_test

import (
	"context"
	"testing"

	"github.com/seinshah/civic/internal/pkg/loader"
	"github.com/stretchr/testify/require"
)

func TestRemoteLoader_Load(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		expectError
		name   string
		path   string
		getCtx func(t *testing.T) (context.Context, context.CancelFunc)
	}{
		{
			name: "success",
			path: "https://google.com",
		},
		{
			name: "timed-out-context",
			path: "https://google.com",
			getCtx: func(t *testing.T) (context.Context, context.CancelFunc) {
				t.Helper()

				return context.WithTimeout(context.Background(), 0)
			},
			expectError: expectError{
				hasError: true,
				err:      context.DeadlineExceeded,
			},
		},
		{
			name: "not-found",
			path: "https://google.com/notfound",
			expectError: expectError{
				hasError: true,
				err:      loader.ErrInvalidRemotePath,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := loader.NewRemoteLoader(tc.path)

			var (
				ctx    context.Context
				cancel context.CancelFunc
			)

			if tc.getCtx != nil {
				ctx, cancel = tc.getCtx(t)

				defer cancel()
			} else {
				ctx = context.Background()
			}

			content, err := l.Load(ctx)

			if tc.hasError {
				require.Error(t, err)

				if tc.err != nil {
					require.ErrorIs(t, err, tc.err)
				}

				return
			}

			require.NoError(t, err)
			require.NotNil(t, content)
		})
	}
}
