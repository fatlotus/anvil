package anvil

import (
	"os"
	"strings"
	"testing"
	"time"
)

func compareTrees(at, bt Tree, t *testing.T) {
	la := 0
	lb := 0

	for {
		a, ma := <-at
		b, mb := <-bt

		if ma {
			la += 1
		}
		if mb {
			lb += 1
		}

		if ma != mb {
			t.Fatalf("trees are of differing length: %d vs. %d", la, lb)
		}

		if ma == false {
			return
		}

		if a.Error() != nil {
			t.Fatal(a.Error())
		}

		if a.Name() != b.Name() {
			t.Fatalf("a.Name() = %s != %s = b.Name()", a.Name(), b.Name())
		}

		dt := a.ModTime().Sub(b.ModTime()).Seconds()

		// FIXME: when Go adds lutimes, add support for symlink metadata.
		if a.Mode()&os.ModeSymlink == 0 {
			if !(-4 < dt && dt < 4) {
				t.Fatalf("a.ModTime() = %s != %s = b.ModTime()",
					a.ModTime(), b.ModTime())
			}

			if a.Mode() != b.Mode() {
				t.Fatalf("a.Mode() = %s != %s = b.Mode()",
					a.Mode(), b.Mode())
			}
		}

		if a.Size() != b.Size() {
			t.Fatalf("a.Size() = %d != %d = b.Size()",
				a.Size(), b.Size())
		}
	}
}

func fixtureTree(blobs []Blob) Tree {
	return makeTree(func(r Tree) {
		for _, b := range blobs {
			r <- b
		}
	})
}

func fixtureOverlayChanges() Tree {
	return fixtureTree([]Blob{
		&memBlob{
			name:     "a/b",
			mode:     os.FileMode(0440),
			size:     17,
			contents: strings.NewReader("new contents of b"),
			modtime:  time.Date(2001, 3, 4, 5, 3, 5, 3, time.UTC),
			source:   "overlay fixture",
		},
		&memBlob{
			name:     "a/c",
			mode:     os.FileMode(0444) | os.ModeSymlink,
			size:     1,
			contents: strings.NewReader("."),
			modtime:  time.Date(2003, 1, 3, 4, 3, 2, 5, time.UTC),
			source:   "fixture",
		},
	})
}

func fixtureOverlayTree() Tree {
	return fixtureTree([]Blob{
		&memBlob{
			name:     "a/",
			mode:     os.FileMode(0755) | os.ModeDir,
			modtime:  time.Date(2000, 3, 12, 1, 2, 3, 4, time.UTC),
			source:   "overlay fixture",
			contents: strings.NewReader(""),
		},
		&memBlob{
			name:     "a/b",
			mode:     os.FileMode(0440),
			size:     17,
			contents: strings.NewReader("new contents of b"),
			modtime:  time.Date(2001, 3, 4, 5, 3, 5, 3, time.UTC),
			source:   "overlay fixture",
		},
		&memBlob{
			name:     "a/c",
			mode:     os.FileMode(0444) | os.ModeSymlink,
			size:     1,
			contents: strings.NewReader("."),
			modtime:  time.Date(2003, 1, 3, 4, 3, 2, 5, time.UTC),
			source:   "fixture",
		},
	})
}

func fixtureValidTree() Tree {
	return fixtureTree([]Blob{
		&memBlob{
			name:     "a/",
			mode:     os.FileMode(0755) | os.ModeDir,
			modtime:  time.Date(2000, 3, 12, 1, 2, 3, 4, time.UTC),
			source:   "fixture",
			contents: strings.NewReader(""),
		},
		&memBlob{
			name:     "a/b",
			mode:     os.FileMode(0440),
			size:     13,
			contents: strings.NewReader("contents of b"),
			modtime:  time.Date(1995, 5, 3, 2, 3, 4, 5, time.UTC),
			source:   "fixture",
		},
		&memBlob{
			name:     "a/c",
			mode:     os.FileMode(0444) | os.ModeSymlink,
			size:     1,
			contents: strings.NewReader("b"),
			modtime:  time.Date(2011, 2, 4, 2, 3, 4, 5, time.UTC),
			source:   "fixture",
		},
		&memBlob{
			name:     "a/d",
			mode:     os.FileMode(0440),
			size:     13,
			contents: strings.NewReader("contents of c"),
			modtime:  time.Date(1986, 2, 5, 3, 1, 4, 5, time.UTC),
			source:   "fixture",
		},
	})
}
