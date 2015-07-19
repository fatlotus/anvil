package hephaestus

import (
	"io"
	"encoding/json"
	"archive/zip"
	"path/filepath"
	"os"
)

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

// The FromZip method allows the generation of a tree given the path to a local
// zip archive.
func FromZip(path string) (Tree, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	return Validate(TreeFromZip(reader, path)), nil
}

// The FromManifest method reads a local MANIFEST.json file.
func FromManifest(path string) (Tree, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	
	manifest := new(manifest)
	json.NewDecoder(fp).Decode(&manifest)
	
	trees := make([]Tree, len(manifest.Entries))
	
	for i, entry := range manifest.Entries {
		tree, err := LoadZip(filepath.Join(filepath.Dir(path), entry.URL))
		if err != nil {
			return nil, err
		}
		
		trees[i] = WithPrefix(tree, filepath.Dir(entry.URL))
	}
	
	return OverlayAll(trees), nil
}

// GenerateManifest reads a directory with packages, adding a manifest to it.
func GenerateManifest(root string, w io.Writer) error {
	manifest := &manifest{
		Entries: make([]manifestEntry, 0),
	}
	
	root = filepath.Clean(root)
	
	err := filepath.Walk(root, func(p string, i os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		
		if filepath.Ext(p) == ".zip" {
			manifest.Entries = append(manifest.Entries, manifestEntry{
				URL: p[len(root) + 1:],
			})
		}
		
		return nil
	})
	
	if err == nil {
		err = json.NewEncoder(w).Encode(&manifest)
	}
	
	return err
}