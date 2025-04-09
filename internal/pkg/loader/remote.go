package loader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
)

type RemoteLoader struct {
	path    string
	content []byte
	mu      sync.RWMutex
}

var _ Loader = (*RemoteLoader)(nil)

func NewRemoteLoader(path string) *RemoteLoader {
	return &RemoteLoader{path: path}
}

func (r *RemoteLoader) Load(ctx context.Context) ([]byte, error) {
	r.mu.RLock()

	if r.content != nil {
		r.mu.RUnlock()

		return r.content, nil
	}

	r.mu.RUnlock()
	r.mu.Lock()
	defer r.mu.Unlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.path, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.Join(ErrInvalidRemotePath, err), r.path)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.Join(ErrInvalidRemotePath, err), r.path)
	}

	if resp == nil {
		return nil, fmt.Errorf("%w: invalid empty response", ErrInvalidRemotePath)
	}

	defer func() {
		if errD := resp.Body.Close(); errD != nil {
			slog.Warn("failed to close the response body", "error", errD, "path", r.path)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s (%s)", ErrInvalidRemotePath, r.path, resp.Status)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.Join(ErrInvalidRemotePath, err), r.path)
	}

	r.content = content

	return r.content, nil
}
