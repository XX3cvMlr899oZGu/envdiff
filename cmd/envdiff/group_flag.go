package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/user/envdiff/internal/envgroup"
)

var (
	groupSeparator  string
	groupStripPrefix bool
)

func registerGroupFlags() {
	flag.StringVar(&groupSeparator, "group-sep", "_", "separator used to split key prefix for grouping")
	flag.BoolVar(&groupStripPrefix, "group-strip", false, "strip prefix from keys in group output")
}

func buildGroupOptions() envgroup.Options {
	return envgroup.Options{
		Separator:   groupSeparator,
		StripPrefix: groupStripPrefix,
	}
}

func printGroups(env map[string]string, opts envgroup.Options) {
	groups := envgroup.Apply(env, opts)
	for _, g := range groups {
		name := g.Name
		if name == "" {
			name = "(no prefix)"
		}
		fmt.Printf("[%s]\n", name)
		keys := make([]string, 0, len(g.Keys))
		for k := range g.Keys {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("  %s=%s\n", k, g.Keys[k])
		}
	}
}
