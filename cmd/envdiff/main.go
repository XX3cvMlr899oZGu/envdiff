package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/report"
)

func main() {
	formatFlag := flag.String("format", "text", "output format: text or json")
	showEqual := flag.Bool("show-equal", false, "include keys with matching values in output")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [flags] <file-a> <file-b>\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	fileA, fileB := args[0], args[1]

	mapA, err := parser.ParseFile(fileA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", fileA, err)
		os.Exit(1)
	}

	mapB, err := parser.ParseFile(fileB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", fileB, err)
		os.Exit(1)
	}

	results := diff.Compare(mapA, mapB)

	opts := report.Options{
		Format:    report.Format(*formatFlag),
		ShowEqual: *showEqual,
	}

	if err := report.Write(os.Stdout, results, fileA, fileB, opts); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(1)
	}

	for _, r := range results {
		if r.Status != diff.StatusEqual {
			os.Exit(1)
		}
	}
}
