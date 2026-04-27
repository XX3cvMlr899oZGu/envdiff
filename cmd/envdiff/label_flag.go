package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/envdiff/internal/envdiff/envdiff/envlabel"
	"github.com/yourusername/envdiff/internal/loader"
)

var (
	labelFile      string
	labelKVs       stringSliceFlag
	labelOverwrite bool
)

func registerLabelFlags(fs *flag.FlagSet) {
	fs.StringVar(&labelFile, "label-file", "", "path to .env file to label")
	fs.Var(&labelKVs, "label", "label in key=value format (repeatable)")
	fs.BoolVar(&labelOverwrite, "label-overwrite", false, "allow overwriting existing labels")
}

func runLabelCommand() {
	if labelFile == "" {
		fmt.Fprintln(os.Stderr, "envdiff label: --label-file is required")
		os.Exit(1)
	}

	env, err := loader.LoadFile(labelFile, loader.Options{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "envdiff label: %v\n", err)
		os.Exit(1)
	}

	labels := make(map[string]string, len(labelKVs))
	for _, kv := range labelKVs {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "envdiff label: invalid label %q, expected key=value\n", kv)
			os.Exit(1)
		}
		labels[parts[0]] = parts[1]
	}

	opts := envlabel.DefaultOptions()
	opts.Labels = labels
	opts.OverwriteExisting = labelOverwrite

	le, err := envlabel.Apply(env, nil, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "envdiff label: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(envlabel.FormatText(le))
}

// stringSliceFlag is a reusable multi-value flag type (defined here if not already present).
type stringSliceFlag []string

func (s *stringSliceFlag) String() string  { return strings.Join(*s, ",") }
func (s *stringSliceFlag) Set(v string) error { *s = append(*s, v); return nil }
