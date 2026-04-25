package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/envdiff/envdiff/envdrift"
	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/snapshot"
)

var (
	driftSnapshotPath   string
	driftLivePath       string
	driftShowUnchanged  bool
)

func registerDriftFlags(fs *flag.FlagSet) {
	fs.StringVar(&driftSnapshotPath, "snapshot", "", "path to snapshot JSON file")
	fs.StringVar(&driftLivePath, "live", "", "path to live .env file")
	fs.BoolVar(&driftShowUnchanged, "show-unchanged", false, "include unchanged keys in output")
}

func runDriftCommand() {
	if driftSnapshotPath == "" || driftLivePath == "" {
		fmt.Fprintln(os.Stderr, "drift: --snapshot and --live are required")
		os.Exit(1)
	}

	snap, err := snapshot.Load(driftSnapshotPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "drift: load snapshot: %v\n", err)
		os.Exit(1)
	}

	live, err := loader.LoadFile(driftLivePath, loader.Options{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "drift: load live env: %v\n", err)
		os.Exit(1)
	}

	opts := envdrift.Options{IncludeUnchanged: driftShowUnchanged}
	entries := envdrift.Detect(snap, live, opts)
	fmt.Print(envdrift.FormatText(entries))

	if envdrift.HasDrift(entries) {
		os.Exit(2)
	}
}
