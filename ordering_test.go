package anvil

import (
	"sort"
	"testing"
)

func fixturePathsOutOfOrder() []string {
	return []string{
		"d",
		"\u60A8\u597D",
		"a/b",
		"a-aa",
		"a/",
		"a-cc",
	}
}

func fixturePathsInOrder() []string {
	return []string{
		"a/",
		"a/b",
		"a-aa",
		"a-cc",
		"d",
		"\u60A8\u597D",
	}
}

func pathsToBlobs(a []string) []Blob {
	r := make([]Blob, len(a))

	for i, path := range a {
		r[i] = &memBlob{
			name: path,
		}
	}

	return r
}

func TestPathOrdering(t *testing.T) {
	a := sortPaths(fixturePathsOutOfOrder())
	sort.Sort(a)

	b := fixturePathsInOrder()

	for i := range a {
		if a[i] != b[i] {
			t.Errorf("a[%d] = %s != %s = b[%d]", i, a[i], b[i], i)
		}
	}
}

func TestBlobOrdering(t *testing.T) {
	a := sortBlobs(pathsToBlobs(fixturePathsOutOfOrder()))
	sort.Sort(a)

	b := pathsToBlobs(fixturePathsInOrder())

	for i := range a {
		if a[i].Name() != b[i].Name() {
			t.Errorf("a[%d] = %s != %s = b[%d]", i, a[i].Name(), b[i].Name(), i)
		}
	}
}
