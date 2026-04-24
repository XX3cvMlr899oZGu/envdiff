// Package envtag assigns named labels (tags) to env keys and groups
// environment entries by those labels.
//
// Tags are defined as a list of Tag values, each carrying a name and the
// set of keys that belong to it. A single key may belong to more than one
// tag. Keys that match no tag are collected under the reserved group name
// "_untagged".
//
// Typical usage:
//
//	opts := envtag.Options{
//	    Tags: []envtag.Tag{
//	        {Name: "database", Keys: []string{"DB_HOST", "DB_PORT"}},
//	        {Name: "secrets",  Keys: []string{"DB_PASSWORD", "API_KEY"}},
//	    },
//	}
//	result, err := envtag.Apply(env, opts)
//	dbKeys := envtag.KeysForTag(result, "database")
package envtag
