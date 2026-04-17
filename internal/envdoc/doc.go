// Package envdoc parses .env files and extracts documentation from
// inline and preceding comments.
//
// It supports two comment styles:
//
//	# This is a preceding comment
//	KEY=value
//
//	KEY=value # This is an inline comment
//
// Preceding comments are associated with the immediately following key.
// A blank line between a comment and a key resets the association.
// Inline comments take precedence over preceding comments.
//
// Use Parse to extract entries from a file, and ToMap to get a
// key-to-comment lookup map.
package envdoc
