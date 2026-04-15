package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envdiff/internal/filter"
	"github.com/user/envdiff/internal/parser"
)

// Options controls how env files are loaded and filtered.
type Options struct {
	Prefix  string
	Exclude []string
	Regex   string
}

// LoadFile reads and parses a .env file, then applies any filter options.
// It returns the resulting key-value map or an error.
func LoadFile(path string, opts Options) (map[string]string, error) {
	path = filepath.Clean(path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	env, err := parser.ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	filtered, err := filter.ApplyToMap(env, filter.Options{
		Prefix:  opts.Prefix,
		Exclude: opts.Exclude,
		Regex:   opts.Regex,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to apply filters: %w", err)
	}

	return filtered, nil
}

// LoadFiles loads multiple .env files and returns a slice of named maps.
// The name for each entry is the file's base name.
func LoadFiles(paths []string, opts Options) ([]NamedEnv, error) {
	results := make([]NamedEnv, 0, len(paths))
	for _, p := range paths {
		env, err := LoadFile(p, opts)
		if err != nil {
			return nil, err
		}
		results = append(results, NamedEnv{
			Name: filepath.Base(p),
			Env:  env,
		})
	}
	return results, nil
}

// NamedEnv pairs a display name with its parsed environment map.
type NamedEnv struct {
	Name string
	Env  map[string]string
}
