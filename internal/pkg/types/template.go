package types

import (
	"html/template"
	"regexp"
	"strings"
)

//go:generate go tool go-enum --nocase --names

var reDomain = regexp.MustCompile(`^(?:http://|https://)?(?:www\.)?([\w-]+)\.\w+(.*)$`)

const TemplateRegistryPath = "https://raw.githubusercontent.com/seinshah/civic/refs/heads/main/templates"

// SocialMediaPlatform is the data representing the name of a social media.
// ENUM(facebook, github, gitlab, linkedin, mastodon, reddit, stackoverflow, x-twitter, youtube, other).
type SocialMediaPlatform string

type SocialMediaLink struct {
	Name             SocialMediaPlatform
	Link             string
	DetectedUsername string
}

type TemplateData struct {
	Schema *Schema
}

func (sbc *SchemaBioContact) ParsedSocials() []SocialMediaLink {
	if sbc == nil || len(sbc.Socials) == 0 {
		return nil
	}

	mediaLinks := make([]SocialMediaLink, 0, len(sbc.Socials))

	for _, social := range sbc.Socials {
		linkParts := reDomain.FindStringSubmatch(social)

		var (
			domain    string
			remainder string
		)

		if len(linkParts) == 3 { //nolint:mnd
			domain = linkParts[1]
			remainder = linkParts[2]
		}

		parsedDomain, err := ParseSocialMediaPlatform(domain)
		if err != nil {
			if strings.ToLower(domain) == "twitter" || strings.ToLower(domain) == "x" {
				parsedDomain = SocialMediaPlatformXTwitter
			} else {
				parsedDomain = SocialMediaPlatformOther
			}
		}

		link := SocialMediaLink{
			Name:             parsedDomain,
			Link:             social,
			DetectedUsername: social,
		}

		remainderPattern := socialDomainUsernameRegexPatter(parsedDomain)

		if remainderPattern != "" {
			remainderRE := regexp.MustCompile(remainderPattern)
			remainderParts := remainderRE.FindStringSubmatch(remainder)

			if len(remainderParts) == 2 { //nolint:mnd
				link.DetectedUsername = remainderParts[1]
			}
		}

		mediaLinks = append(mediaLinks, link)
	}

	return mediaLinks
}

func UnescapeHTML(s string) template.HTML {
	return template.HTML(s) //nolint:gosec
}

func socialDomainUsernameRegexPatter(domain SocialMediaPlatform) string {
	var pattern string

	switch domain {
	case SocialMediaPlatformFacebook, SocialMediaPlatformGithub, SocialMediaPlatformGitlab,
		SocialMediaPlatformXTwitter:
		pattern = `^/([\w-\.]+).*$`

	case SocialMediaPlatformReddit:
		pattern = `^/user/([\w-\.]+).*$`

	case SocialMediaPlatformStackoverflow:
		pattern = `^/users/\d+/([\w-\.]+).*$`

	case SocialMediaPlatformLinkedin:
		pattern = `^/([\w-]+/[\w-\.]+).*$`

	case SocialMediaPlatformMastodon, SocialMediaPlatformYoutube:
		pattern = `^/@([\w-\.]+).*$`
	}

	return pattern
}
