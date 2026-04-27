// Package envlabel provides utilities for attaching and managing string labels
// on environment maps.
//
// Labels are arbitrary key/value metadata (e.g. env=production, team=platform)
// that travel alongside an env map without polluting its keys. They are useful
// for routing, auditing, and display purposes in multi-environment pipelines.
//
// Basic usage:
//
//	opts := envlabel.DefaultOptions()
//	opts.Labels = map[string]string{"env": "production"}
//	le, err := envlabel.Apply(env, existingLabels, opts)
//	fmt.Print(envlabel.FormatText(le))
package envlabel
