package version_test

import (
	"testing"

	"github.com/seinshah/civic/internal/pkg/version"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		version string
		want    *version.Semantic
		wantErr bool
	}{
		{
			name:    "valid full version with v",
			version: "v1.2.3",
			want:    version.New(1, 2, 3),
		},
		{
			name:    "valid full version without v",
			version: "1.2.3",
			want:    version.New(1, 2, 3),
		},
		{
			name:    "valid major version with v",
			version: "v1",
			want:    version.New(1, 0, 0),
		},
		{
			name:    "valid major version without v",
			version: "1",
			want:    version.New(1, 0, 0),
		},
		{
			name:    "valid major and minor version with v",
			version: "v1.2",
			want:    version.New(1, 2, 0),
		},
		{
			name:    "valid major and minor version without v",
			version: "1.2",
			want:    version.New(1, 2, 0),
		},
		{
			name:    "valid with arbitrary build or pre-release",
			version: "v1.2.3-alpha+build",
			want:    version.New(1, 2, 3),
		},
		{
			name:    "invalid version",
			version: "invalid",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := version.Parse(tc.version)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, got)

				return
			}

			require.NoError(t, err)
			require.True(t, tc.want.Equal(got))
		})
	}
}
