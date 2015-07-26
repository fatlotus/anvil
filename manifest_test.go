package anvil

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Dump the given Stream onto the particular location on disk.
func dumpZip(tree Stream, path string, t *testing.T) {
	fp, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	if err := tree.ToZip(zip.NewWriter(fp)); err != nil {
		t.Fatal(err)
	}
}

// A simple manifest fixture.
func fixtureManifest() Stream {
	return fixtureStream([]Blob{
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
		&memBlob{
			name:     "a/d",
			mode:     os.FileMode(0440),
			size:     13,
			contents: strings.NewReader("contents of c"),
			modtime:  time.Date(1986, 2, 5, 3, 1, 4, 5, time.UTC),
			source:   "fixture",
		},
		&memBlob{
			name:     "d/",
			mode:     os.FileMode(0755) | os.ModeDir,
			modtime:  time.Now(),
			source:   "fixture",
			contents: strings.NewReader(""),
		},
		&memBlob{
			name:     "d/a/",
			mode:     os.FileMode(0755) | os.ModeDir,
			modtime:  time.Date(2000, 3, 12, 1, 2, 3, 4, time.UTC),
			source:   "fixture",
			contents: strings.NewReader(""),
		},
		&memBlob{
			name:     "d/a/b",
			mode:     os.FileMode(0440),
			size:     13,
			contents: strings.NewReader("contents of b"),
			modtime:  time.Date(1995, 5, 3, 2, 3, 4, 5, time.UTC),
			source:   "fixture",
		},
		&memBlob{
			name:     "d/a/c",
			mode:     os.FileMode(0444) | os.ModeSymlink,
			size:     1,
			contents: strings.NewReader("b"),
			modtime:  time.Date(2011, 2, 4, 2, 3, 4, 5, time.UTC),
			source:   "fixture",
		},
		&memBlob{
			name:     "d/a/d",
			mode:     os.FileMode(0440),
			size:     13,
			contents: strings.NewReader("contents of c"),
			modtime:  time.Date(1986, 2, 5, 3, 1, 4, 5, time.UTC),
			source:   "fixture",
		},
	})
}

// Verifies the integrity of file:// and http:// images.
func TestManifests(t *testing.T) {
	// Create a test image directory.
	tmp, err := ioutil.TempDir("", "heph-man-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	if err := os.Mkdir(filepath.Join(tmp, "d"), os.FileMode(0755)); err != nil {
		t.Fatal(err)
	}

	// Add several overlays to the image.
	dumpZip(fixtureValidStream(), filepath.Join(tmp, "A.zip"), t)
	dumpZip(fixtureValidStream(), filepath.Join(tmp, "d", "C.zip"), t)
	dumpZip(fixtureOverlayStream(), filepath.Join(tmp, "B.zip"), t)

	if err := RecomputeManifest(tmp); err != nil {
		t.Fatal(err)
	}

	// Compute the result of composing each archive.
	img, err := FromImage(url.URL{Path: tmp, Scheme: "file"})
	if err != nil {
		fmt.Printf("tmp: %s; err: %s\n", tmp, err)
		time.Sleep(60 * time.Second)
	}

	// Compute the result over HTTP.
	server := httptest.NewServer(http.FileServer(http.Dir(tmp)))
	defer server.Close()

	httpimg, err := FromImage(url.URL{Path: tmp, Scheme: "file"})
	if err != nil {
		fmt.Printf("tmp: %s; err: %s\n", tmp, err)
		time.Sleep(60 * time.Second)
	}

	// Compare against the refernce.
	compareStreams(img, fixtureManifest(), t)
	compareStreams(httpimg, fixtureManifest(), t)
}
