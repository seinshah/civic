package types_test

import (
	"testing"

	"github.com/seinshah/civic/internal/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestSchemaBioContact_ParsedSocials(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		data     *types.SchemaBioContact
		expected []types.SocialMediaLink
	}{
		{
			name: "no bio contact",
		},
		{
			name: "no bio contact socials",
			data: &types.SchemaBioContact{},
		},
		{
			name: "mixed items",
			data: &types.SchemaBioContact{
				Socials: []string{
					"https://www.facebook.com/someuser",
					"https://github.com/someuser",
					"https://gitlab.com/someuser",
					"https://www.linkedin.com/in/someuser/",
					"https://mastodon.social/@someuser",
					"https://www.reddit.com/user/someuser/",
					"https://stackoverflow.com/users/123456/someuser",
					"https://x.com/someuser",
					"https://www.youtube.com/@someuser",
					"https://othersocial.com/something/username=someuser",
				},
			},
			expected: []types.SocialMediaLink{
				{
					Name: types.SocialMediaPlatformFacebook,
				},
				{
					Name: types.SocialMediaPlatformGithub,
				},
				{
					Name: types.SocialMediaPlatformGitlab,
				},
				{
					Name:             types.SocialMediaPlatformLinkedin,
					DetectedUsername: "in/someuser",
				},
				{
					Name: types.SocialMediaPlatformMastodon,
				},
				{
					Name: types.SocialMediaPlatformReddit,
				},
				{
					Name: types.SocialMediaPlatformStackoverflow,
				},
				{
					Name: types.SocialMediaPlatformXTwitter,
				},
				{
					Name: types.SocialMediaPlatformYoutube,
				},
				{
					Name:             types.SocialMediaPlatformOther,
					DetectedUsername: "https://othersocial.com/something/username=someuser",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				t.Parallel()

				if len(tc.expected) > 0 {
					for i := range tc.expected {
						tc.expected[i].Link = tc.data.Socials[i]

						if tc.expected[i].DetectedUsername == "" {
							tc.expected[i].DetectedUsername = "someuser"
						}
					}
				}

				actual := tc.data.ParsedSocials()

				require.Equal(t, tc.expected, actual)
			},
		)
	}
}
