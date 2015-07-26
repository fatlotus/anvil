package anvil

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/DHowett/ranger"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// A Resource is a blob that mostly acts like an os.File.
type resource interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

type rangerResource struct {
	name string
	*ranger.Reader
}

func (r *rangerResource) Stat() (os.FileInfo, error) { return r, nil }
func (r *rangerResource) Name() string               { return r.name }
func (r *rangerResource) IsDir() bool                { return false }
func (r *rangerResource) ModTime() time.Time         { return time.Now() }
func (r *rangerResource) Mode() os.FileMode          { return 0 }
func (r *rangerResource) Size() int64                { return r.Length() }
func (r *rangerResource) Sys() interface{}           { return nil }

// Opens the given URL, returning a Resource.
func openURL(url *url.URL) (resource, error) {
	switch url.Scheme {
	case "http", "https":

		// Resolve the URL as a remote HTTP/HTTPS stream.
		reader, err := ranger.NewReader(&ranger.HTTPRanger{URL: url})
		if err != nil {
			return nil, err
		}

		return &rangerResource{url.String(), reader}, nil

	case "file", "":

		// Resolve the URL as a local file.
		return os.Open(url.Path)

	default:
		return nil, fmt.Errorf("Unknown URL scheme: %s", url)
	}
}

// The FromImage loader resolves an image as a manifest file.
func FromImage(base url.URL) (Stream, error) {

	// Read the manifest as a file.
	base.Path += "/"

	res, err := openURL(
		base.ResolveReference(&url.URL{Path: "./MANIFEST.json"}))
	if err != nil {
		return nil, err
	}

	// Parse the entry as a manifest file.
	var m manifest
	json.NewDecoder(res).Decode(&m)

	// Create a Stream object for each manifest entry.
	trees := make([]Stream, len(m.Entries))

	for i, entry := range m.Entries {
		parsed, err := url.Parse(entry.URL)
		if err != nil {
			return nil, err
		}

		res, err := openURL(base.ResolveReference(parsed))
		if err != nil {
			return nil, err
		}

		stat, err := res.Stat()
		if err != nil {
			return nil, err
		}

		reader, err := zip.NewReader(res, stat.Size())
		if err != nil {
			return nil, err
		}

		t := FromZip(reader, stat.Name()).WithPrefix(filepath.Dir(entry.URL))
		trees[i] = t
	}

	return OverlayAll(trees), nil
}

// The manifestEntry type represents a subpackage of a manifest file of the
// given type.
type manifestEntry struct {
	URL string `json:"url"`
}

// A manifest file contains a list of all patches required to construct a given
// machine state.
type manifest struct {
	Entries []manifestEntry `json:"entries"`
}

// RecomputeManifest reads a directory with packages, adding a manifest to it.
func RecomputeManifest(root string) error {
	manifest := &manifest{
		Entries: make([]manifestEntry, 0),
	}

	root = filepath.Clean(root)

	err := filepath.Walk(root, func(p string, i os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		if filepath.Ext(p) == ".zip" {
			rel, err := filepath.Rel(root, p)
			if err != nil {
				return err
			}

			manifest.Entries = append(manifest.Entries, manifestEntry{
				URL: rel,
			})
		}

		return nil
	})

	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(root, "MANIFEST.json"))
	if err != nil {
		return err
	}

	if err := json.NewEncoder(file).Encode(&manifest); err != nil {
		return err
	}

	return nil
}
