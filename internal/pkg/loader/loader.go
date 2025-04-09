package loader

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"strings"
	"syscall"
)

type Loader interface {
	Load(ctx context.Context) ([]byte, error)
}

type GeneralLoader struct {
	loader Loader
}

var _ Loader = (*GeneralLoader)(nil)

var (
	ErrInvalidPath       = errors.New("invalid local or remote path")
	ErrInvalidLocalPath  = errors.New("failed to load the local file")
	ErrInvalidRemotePath = errors.New("failed to load the remote file")
	ErrNoFileToLoad      = errors.New("specific loader is not properly configured")
)

func NewGeneralLoader(filePath string) (*GeneralLoader, error) {
	var generic GeneralLoader

	if isRemotePath(filePath) {
		generic.loader = NewRemoteLoader(filePath)
	} else if isLocalPath(filePath) {
		generic.loader = NewLocalLoader(filePath)
	}

	if generic.loader == nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidPath, filePath)
	}

	return &generic, nil
}

func (l *GeneralLoader) Load(ctx context.Context) ([]byte, error) {
	if l == nil || l.loader == nil {
		return nil, ErrNoFileToLoad
	}

	return l.loader.Load(ctx)
}

// Loader simply return the concrete underlyin loader that GeneralLoader has detected
// and will be using going forward.
//
//nolint:ireturn
func (l *GeneralLoader) Loader() Loader {
	return l.loader
}

func isLocalPath(path string) bool {
	// We make sure it isn't a directory.
	if strings.HasSuffix(path, string(os.PathSeparator)) {
		return false
	}

	if _, err := os.Stat(path); err != nil {
		var pathErr *fs.PathError

		if errors.As(err, &pathErr) {
			if errors.Is(pathErr.Err, syscall.EINVAL) {
				// It's definitely an invalid character in the filepath.
				return false
			}
		}
	}

	// It could be a permission error, a does-not-exist error, etc.
	// Out-of-scope for this validation, though.
	return true
}

func isRemotePath(path string) bool {
	parsed, err := url.Parse(path)
	if err != nil {
		return false
	}

	return parsed.Scheme == "http" || parsed.Scheme == "https"
}
