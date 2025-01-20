package version

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/seinshah/civic/internal/pkg/types"
)

var (
	errInvalidVersion = errors.New("invalid version")

	errRequestFailed   = errors.New("failed to get the latest version")
	errParsingResponse = errors.New("failed to parse the response")
	errNoLatestVersion = errors.New("no latest version found")
)

const requestTimeout = 5 * time.Second

type Semantic struct {
	major int
	minor int
	patch int
}

// regex that parses a semver to its major, minor, and patch versions.
// appearance of v, minor, path, and pre-release/build is optional.
var semverRE = regexp.MustCompile(`^v?(\d+)(?:\.(\d+)(?:\.(\d+))?)?.*$`)

func New(major, minor, patch int) *Semantic {
	return &Semantic{
		major: major,
		minor: minor,
		patch: patch,
	}
}

// Parse parses a semantic version and detects major, minor, and patch versions.
// If the provided version is Latest, we will attempt to detect the repository's latest version
// using GitHub releases.
func Parse(version string) (*Semantic, error) {
	matches := semverRE.FindStringSubmatch(version)

	// nolint: mnd
	if len(matches) < 4 {
		return nil, errInvalidVersion
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, errInvalidVersion
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		minor = 0
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		patch = 0
	}

	return &Semantic{
		major: major,
		minor: minor,
		patch: patch,
	}, nil
}

// ParseFromGithub attempts to get the latest version of the repository from GitHub releases.
func ParseFromGithub(ctx context.Context) (*Semantic, error) {
	latest, err := getLatestFromGithub(ctx)
	if err != nil {
		return nil, err
	}

	return Parse(latest)
}

func (s *Semantic) String() string {
	return fmt.Sprintf("v%d.%d.%d", s.major, s.minor, s.patch)
}

func (s *Semantic) Equal(s2 *Semantic) bool {
	return s.major == s2.major && s.minor == s2.minor && s.patch == s2.patch
}

func (s *Semantic) GreaterThan(cmp *Semantic) bool {
	if s.major > cmp.major {
		return true
	}

	if s.major == cmp.major && s.minor > cmp.minor {
		return true
	}

	if s.major == cmp.major && s.minor == cmp.minor && s.patch > cmp.patch {
		return true
	}

	return false
}

func (s *Semantic) Major() int {
	return s.major
}

func getLatestFromGithub(ctx context.Context) (string, error) {
	newCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(
		newCtx,
		http.MethodGet,
		fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/releases?per_page=1",
			types.DefaultAppOwner,
			types.DefaultAppName,
		),
		nil,
	)
	if err != nil {
		return "", errors.Join(errRequestFailed, err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28") // nolint: canonicalheader

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Join(errRequestFailed, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: %s", errRequestFailed, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Join(errParsingResponse, err)
	}

	data := make([]map[string]any, 0)

	if err = json.Unmarshal(body, &data); err != nil {
		return "", errors.Join(errParsingResponse, err)
	}

	if len(data) == 0 {
		return "", errNoLatestVersion
	}

	tagName, ok := data[0]["tag_name"].(string)
	if !ok {
		return "", errNoLatestVersion
	}

	return tagName, nil
}
