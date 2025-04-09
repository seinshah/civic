package chrome_test

import (
	"context"
	"testing"
	"time"

	"github.com/seinshah/civic/internal/pkg/output/pdf/chrome"
	"github.com/seinshah/civic/internal/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestEngine_Generate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		options     []chrome.Option
		getCtx      func(t *testing.T) (context.Context, func())
		expectError bool
		err         error
	}{
		{
			name: "success-without-options",
		},
		{
			name: "success-with-options",
			options: []chrome.Option{
				chrome.WithPageSize(types.PageSizeA4),
				chrome.WithPageMargin(
					types.PageMargin{
						Top:    2,
						Right:  2,
						Bottom: 2,
						Left:   2,
					},
				),
			},
		},
		{
			name: "timed-out-context",
			getCtx: func(t *testing.T) (context.Context, func()) {
				t.Helper()

				return context.WithTimeout(t.Context(), time.Duration(0))
			},
			expectError: true,
			err:         context.DeadlineExceeded,
		},
		{
			name: "invalid-page-size",
			options: []chrome.Option{
				chrome.WithPageSize("invalid"),
			},
			expectError: true,
			err:         types.ErrInvalidPageSize,
		},
		{
			name: "invalid-page-margin",
			options: []chrome.Option{
				chrome.WithPageMargin(
					types.PageMargin{
						Top:    3,
						Right:  3,
						Bottom: 3,
						Left:   3,
					},
				),
			},
			expectError: true,
			err:         types.ErrInvalidPageMargin,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				t.Parallel()

				var (
					ctx    = t.Context()
					cancel func()
				)

				if tc.getCtx != nil {
					ctx, cancel = tc.getCtx(t)
					defer cancel()
				}

				engine := chrome.NewHeadless(tc.options...)
				output, err := engine.Generate(ctx, []byte("<p>test</p>"))

				if tc.expectError {
					require.Error(t, err)
					require.Nil(t, output)

					return
				}

				require.NoError(t, err)
				require.NotEmpty(t, output)
			},
		)
	}
}
