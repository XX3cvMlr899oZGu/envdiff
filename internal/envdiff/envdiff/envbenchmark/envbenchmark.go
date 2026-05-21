// Package envbenchmark measures parse and diff performance across env maps.
package envbenchmark

import (
	"fmt"
	"sort"
	"time"
)

// Result holds the outcome of a single benchmark run.
type Result struct {
	Label    string
	Duration time.Duration
	KeyCount int
}

// Options configures benchmark behaviour.
type Options struct {
	// Runs is the number of times each operation is executed; defaults to 1.
	Runs int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Runs: 1}
}

// Run executes fn the configured number of times and returns the mean duration.
func Run(label string, env map[string]string, fn func(map[string]string) error, opts Options) (Result, error) {
	if opts.Runs < 1 {
		opts.Runs = 1
	}

	var total time.Duration
	for i := 0; i < opts.Runs; i++ {
		start := time.Now()
		if err := fn(env); err != nil {
			return Result{}, fmt.Errorf("envbenchmark: run %d failed: %w", i+1, err)
		}
		total += time.Since(start)
	}

	return Result{
		Label:    label,
		Duration: total / time.Duration(opts.Runs),
		KeyCount: len(env),
	}, nil
}

// FormatText renders a slice of Results as a human-readable table.
func FormatText(results []Result) string {
	if len(results) == 0 {
		return "no benchmark results\n"
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Duration < results[j].Duration
	})

	out := fmt.Sprintf("%-30s %12s %8s\n", "label", "mean_ns", "keys")
	out += fmt.Sprintf("%-30s %12s %8s\n", "-----", "-------", "----")
	for _, r := range results {
		out += fmt.Sprintf("%-30s %12d %8d\n", r.Label, r.Duration.Nanoseconds(), r.KeyCount)
	}
	return out
}
