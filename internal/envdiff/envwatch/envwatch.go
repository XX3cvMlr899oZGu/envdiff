// Package envwatch provides functionality to watch .env files for changes
// and report diffs when the file content changes on disk.
package envwatch

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/envdiff/internal/parser"
)

// Event represents a change detected in a watched .env file.
type Event struct {
	Path    string
	OldEnv  map[string]string
	NewEnv  map[string]string
	OccurredAt time.Time
}

// Options controls the behaviour of the watcher.
type Options struct {
	// PollInterval is how often the file is checked for changes.
	PollInterval time.Duration
}

// DefaultOptions returns sensible defaults for watching.
func DefaultOptions() Options {
	return Options{
		PollInterval: 2 * time.Second,
	}
}

// Watch polls the given .env file at the configured interval and sends an
// Event on the returned channel whenever the file content changes.
// The caller must close the done channel to stop watching.
func Watch(path string, opts Options, done <-chan struct{}) (<-chan Event, error) {
	env, err := parser.ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("envwatch: initial parse of %q: %w", path, err)
	}

	lastHash, err := fileHash(path)
	if err != nil {
		return nil, fmt.Errorf("envwatch: initial hash of %q: %w", path, err)
	}

	ch := make(chan Event, 1)

	go func() {
		defer close(ch)
		ticker := time.NewTicker(opts.PollInterval)
		defer ticker.Stop()
		current := env
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				h, err := fileHash(path)
				if err != nil || h == lastHash {
					continue
				}
				newEnv, err := parser.ParseFile(path)
				if err != nil {
					continue
				}
				ch <- Event{
					Path:       path,
					OldEnv:     current,
					NewEnv:     newEnv,
					OccurredAt: time.Now(),
				}
				lastHash = h
				current = newEnv
			}
		}
	}()

	return ch, nil
}

func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
