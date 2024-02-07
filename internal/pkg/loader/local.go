package loader

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
)

type LocalLoader struct {
	path    string
	content []byte
	mu      sync.RWMutex
}

var _ Loader = (*LocalLoader)(nil)

func NewLocalLoader(path string) *LocalLoader {
	return &LocalLoader{path: path}
}

func (l *LocalLoader) Load(_ context.Context) ([]byte, error) {
	l.mu.RLock()

	if l.content != nil {
		l.mu.RUnlock()

		return l.content, nil
	}

	l.mu.RUnlock()
	l.mu.Lock()
	defer l.mu.Unlock()

	if stat, err := os.Stat(l.path); err != nil {
		return nil, fmt.Errorf("%w: %s", errors.Join(ErrInvalidLocalPath, err), l.path)
	} else if stat.IsDir() {
		return nil, fmt.Errorf("%w (directory): %s", ErrInvalidLocalPath, l.path)
	}

	content, err := os.ReadFile(l.path)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.Join(ErrInvalidLocalPath, err), l.path)
	}

	l.content = content

	return l.content, nil
}
