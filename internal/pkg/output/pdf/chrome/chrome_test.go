package chrome_test

import (
	"context"
	"testing"
	"time"

	"github.com/seinshah/cvci/internal/pkg/output/pdf/chrome"
	"github.com/seinshah/cvci/internal/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestEngine_Generate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		options     []chrome.Option
		getCtx      func() (context.Context, func())
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
				chrome.WithPageMargin(types.PageMargin{
					Top:    2,
					Right:  2,
					Bottom: 2,
					Left:   2,
				}),
			},
		},
		{
			name: "timed-out-context",
			getCtx: func() (context.Context, func()) {
				return context.WithTimeout(context.Background(), time.Duration(0))
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
				chrome.WithPageMargin(types.PageMargin{
					Top:    3,
					Right:  3,
					Bottom: 3,
					Left:   3,
				}),
			},
			expectError: true,
			err:         types.ErrInvalidPageMargin,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var (
				ctx    context.Context
				cancel func()
			)

			if tc.getCtx != nil {
				ctx, cancel = tc.getCtx()
				defer cancel()
			} else {
				ctx = context.Background()
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
		})
	}
}
