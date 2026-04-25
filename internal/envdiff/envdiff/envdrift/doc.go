// Package envdrift compares a previously saved snapshot of an environment
// against a current (live) environment map and reports configuration drift.
//
// Drift is categorised as one of four statuses:
//
//   - added    – key exists in live but not in the snapshot
//   - removed  – key existed in the snapshot but is absent from live
//   - changed  – key exists in both but its value has been modified
//   - unchanged – key exists in both with identical values
//
// Basic usage:
//
//	snap, _ := snapshot.Load("snapshot.json")
//	live, _ := loader.LoadFile(".env", loader.Options{})
//	entries := envdrift.Detect(snap, live, envdrift.DefaultOptions())
//	if envdrift.HasDrift(entries) {
//		fmt.Print(envdrift.FormatText(entries))
//	}
package envdrift
