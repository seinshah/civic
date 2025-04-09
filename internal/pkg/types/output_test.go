package types_test

import (
	"testing"

	"github.com/seinshah/civic/internal/pkg/types"
	"github.com/stretchr/testify/require"
)

type customString interface {
	IsValid() bool
}

func TestDetectOutputType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		got     func() customString
		isValid bool
		want    any
	}{
		{
			name: "detect pdf",
			got: func() customString {
				return types.DetectFileType[types.OutputType]("path/to/test.pdf")
			},
			isValid: true,
			want:    types.OutputTypePdf,
		},
		{
			name: "detect html",
			got: func() customString {
				return types.DetectFileType[types.OutputType]("path/to/test.html")
			},
			isValid: true,
			want:    types.OutputTypeHtml,
		},
		{
			name: "multi extension",
			got: func() customString {
				return types.DetectFileType[types.OutputType]("x.txt.pdf")
			},
			isValid: true,
			want:    types.OutputTypePdf,
		},
		{
			name: "detect unknown",
			got: func() customString {
				return types.DetectFileType[types.OutputType]("test.txt")
			},
			isValid: false,
		},
		{
			name: "random string",
			got: func() customString {
				return types.DetectFileType[types.OutputType]("doasogfnasigoasignoiasng")
			},
			isValid: false,
		},
		{
			name: "reverse multi extension",
			got: func() customString {
				return types.DetectFileType[types.OutputType]("x.pdf.txt")
			},
			isValid: false,
		},
		{
			name: "detect yaml",
			got: func() customString {
				return types.DetectFileType[types.SchemaType]("path/to/test.yaml")
			},
			isValid: true,
			want:    types.SchemaTypeYaml,
		},
		{
			name: "yml",
			got: func() customString {
				return types.DetectFileType[types.SchemaType]("path/to/test.ex.yml")
			},
			isValid: true,
			want:    types.SchemaTypeYml,
		},
		{
			name: "json",
			got: func() customString {
				return types.DetectFileType[types.SchemaType]("path/to/test.json")
			},
			isValid: false,
		},
		{
			name: "toml",
			got: func() customString {
				return types.DetectFileType[types.SchemaType]("path/to/test.toml")
			},
			isValid: false,
		},
	}
	for _, tc := range tests {
		t.Run(
			tc.name, func(t *testing.T) {
				t.Parallel()

				got := tc.got()

				require.Equal(t, tc.isValid, got.IsValid())

				if tc.isValid {
					require.Equal(t, tc.want, got)
				}
			},
		)
	}
}
