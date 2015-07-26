package anvil

import (
	"os"
	"strings"
	"testing"
	"time"
)

func fixtureStreamWithPrefix() Stream {
	return fixtureStream([]Blob{
		&memBlob{
			name:     "outer/",
			mode:     os.FileMode(0755) | os.ModeDir,
			modtime:  time.Now(),
			contents: strings.NewReader(""),
		},
		&memBlob{
			name:     "outer/inner/",
			mode:     os.FileMode(0755) | os.ModeDir,
			modtime:  time.Now(),
			contents: strings.NewReader(""),
		},
		&memBlob{
			name:     "outer/inner/a/",
			mode:     os.FileMode(0755) | os.ModeDir,
			modtime:  time.Date(2000, 3, 12, 1, 2, 3, 4, time.UTC),
			source:   "fixture",
			contents: strings.NewReader(""),
		},
		&memBlob{
			name:     "outer/inner/a/b",
			mode:     os.FileMode(0440),
			size:     13,
			contents: strings.NewReader("contents of b"),
			modtime:  time.Date(1995, 5, 3, 2, 3, 4, 5, time.UTC),
			source:   "fixture",
		},
		&memBlob{
			name:     "outer/inner/a/c",
			mode:     os.FileMode(0444) | os.ModeSymlink,
			size:     1,
			contents: strings.NewReader("b"),
			modtime:  time.Date(2011, 2, 4, 2, 3, 4, 5, time.UTC),
			source:   "fixture",
		},
		&memBlob{
			name:     "outer/inner/a/d",
			mode:     os.FileMode(0440),
			size:     13,
			contents: strings.NewReader("contents of d"),
			modtime:  time.Date(1986, 2, 5, 3, 1, 4, 5, time.UTC),
			source:   "fixture",
		},
	})
}

func TestWithPrefix(t *testing.T) {
	a := fixtureStreamWithPrefix()
	b := fixtureValidStream().WithPrefix("outer/whoops/.././inner/")

	compareStreams(a, b, t)
}

func TestNoSuperfluous(t *testing.T) {
	compareStreams(fixtureValidStream(), fixtureValidStream().WithPrefix("./"), t)
}

func TestOverlayAll(t *testing.T) {
	result := OverlayAll([]Stream{fixtureValidStream(), fixtureOverlayStream()})
	diff := Difference(fixtureValidStream(), result)

	compareStreams(diff, fixtureOverlayChanges(), t)
}
