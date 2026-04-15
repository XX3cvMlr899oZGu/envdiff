// Package rename provides key-renaming functionality for env maps.
//
// It supports two modes of operation:
//
//  1. Explicit mapping: supply a map[string]string where each entry maps
//     an old key name to its desired new name.
//
//  2. Prefix substitution: supply OldPrefix and NewPrefix in Options to
//     bulk-rename all keys that begin with OldPrefix.
//
// Both modes can be combined; the explicit Map is applied first, and prefix
// substitution is only applied to keys that were not already renamed via the
// explicit map.
//
// Conflicts (two source keys mapping to the same output key) are reported as
// errors so no data is silently overwritten.
//
// Example:
//
//	opts := rename.Options{
//	    Map:       map[string]string{"LEGACY_URL": "DATABASE_URL"},
//	    OldPrefix: "DEV_",
//	    NewPrefix: "PROD_",
//	}
//	res, err := rename.Apply(env, opts)
package rename
