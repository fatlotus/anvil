package anvil

import (
	"strings"
)

// The manglePath comparator compares two string paths componentwise.
func mungePath(a string) string {
	return strings.Replace(a, "/", "\x00", -1)
}

func isPathLessThan(a, b string) bool {
	return mungePath(a) < mungePath(b)
}

func isBlobLessThan(a, b Blob) bool {
	if b == nil {
		return (a != nil)
	} else if a == nil {
		return false
	}

	return isPathLessThan(a.Name(), b.Name())
}

// SortBlobs allows blobs to be sorted componentwise
type sortBlobs []Blob

func (s sortBlobs) Len() int           { return len(s) }
func (s sortBlobs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortBlobs) Less(i, j int) bool { return isBlobLessThan(s[i], s[j]) }

// SortPaths allows strings to be sorted componentwise
type sortPaths []string

func (s sortPaths) Len() int           { return len(s) }
func (s sortPaths) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortPaths) Less(i, j int) bool { return isPathLessThan(s[i], s[j]) }
