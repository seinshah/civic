package types_test

import (
	"testing"

	"github.com/seinshah/cvci/internal/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestDetectOutputType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		isValid bool
		want    types.OutputType
	}{
		{
			name:    "detect pdf",
			path:    "path/to/test.pdf",
			isValid: true,
			want:    types.OutputTypePdf,
		},
		{
			name:    "detect html",
			path:    "path/to/test.html",
			isValid: true,
			want:    types.OutputTypeHtml,
		},
		{
			name:    "multi extension",
			path:    "x.txt.pdf",
			isValid: true,
			want:    types.OutputTypePdf,
		},
		{
			name:    "detect unknown",
			path:    "test.txt",
			isValid: false,
		},
		{
			name:    "random string",
			path:    "doasogfnasigoasignoiasng",
			isValid: false,
		},
		{
			name:    "reverse multi extension",
			path:    "x.pdf.txt",
			isValid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := types.DetectOutputType(tc.path)

			require.Equal(t, tc.isValid, got.IsValid())

			if tc.isValid {
				require.Equal(t, tc.want, got)
			}
		})
	}
}
