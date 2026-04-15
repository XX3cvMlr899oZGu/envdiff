package main

import (
	"flag"
	"strings"

	"github.com/user/envdiff/internal/rename"
)

// keyMapping is a flag.Value that accumulates OLD=NEW pairs.
type keyMapping []string

func (m *keyMapping) String() string { return strings.Join(*m, ",") }
func (m *keyMapping) Set(v string) error {
	*m = append(*m, v)
	return nil
}

var (
	renameKeys   keyMapping
	renameOldPfx string
	renameNewPfx string
)

// registerRenameFlags attaches rename-related flags to the default FlagSet.
func registerRenameFlags() {
	flag.Var(&renameKeys, "rename", "rename a key: OLD=NEW (repeatable)")
	flag.StringVar(&renameOldPfx, "rename-old-prefix", "", "prefix to replace in key names")
	flag.StringVar(&renameNewPfx, "rename-new-prefix", "", "replacement prefix for key names")
}

// buildRenameOptions constructs a rename.Options from parsed CLI flags.
func buildRenameOptions() (rename.Options, error) {
	opts := rename.Options{
		Map:       make(map[string]string),
		OldPrefix: renameOldPfx,
		NewPrefix: renameNewPfx,
	}
	for _, pair := range renameKeys {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return opts, nil
		}
		opts.Map[parts[0]] = parts[1]
	}
	return opts, nil
}
