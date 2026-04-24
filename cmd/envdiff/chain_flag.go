package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envdiff/internal/envchain"
	"github.com/yourorg/envdiff/internal/envtrim"
	"github.com/yourorg/envdiff/internal/filter"
	"github.com/yourorg/envdiff/internal/loader"
)

var (
	chainPrefix  string
	chainTrim    bool
	chainVerbose bool
)

func registerChainFlags() {
	flag.StringVar(&chainPrefix, "chain-prefix", "", "filter keys by prefix before running chain")
	flag.BoolVar(&chainTrim, "chain-trim", false, "trim whitespace from keys and values in chain")
	flag.BoolVar(&chainVerbose, "chain-verbose", false, "print step names before running")
}

func runChainCommand(path string) error {
	env, err := loader.LoadFile(path, loader.Options{})
	if err != nil {
		return fmt.Errorf("load %q: %w", path, err)
	}

	chain := envchain.New()

	if chainPrefix != "" {
		pfx := chainPrefix
		chain.Add("filter-prefix", func(e map[string]string) (map[string]string, error) {
			return filter.ApplyToMap(e, filter.Options{Prefix: pfx}), nil
		})
	}

	if chainTrim {
		chain.Add("trim", func(e map[string]string) (map[string]string, error) {
			return envtrim.Apply(e, envtrim.DefaultOptions()), nil
		})
	}

	if chainVerbose {
		fmt.Fprintf(os.Stderr, "chain steps (%d): %v\n", chain.StepCount(), chain.StepNames())
	}

	result, err := chain.Run(env)
	if err != nil {
		return err
	}

	for k, v := range result {
		fmt.Printf("%s=%s\n", k, v)
	}
	return nil
}
