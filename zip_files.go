package hephaestus

import (
	"archive/zip"
	"sort"
	"fmt"
	"io"
	"strings"
)

const uint32max = 4294967295

type blobSorter []Blob

func (s blobSorter) Len() int { return len(s) }
func (s blobSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s blobSorter) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

// Reads a tree from the given zip stream, marking the blobs as being from the
// given source.
func FromZipFile(r *zip.ReadCloser, source string) Tree {
	return makeTree(func (result Tree) {
		
		blobs := make([]Blob, 0)
		
		for _, f := range r.File {
			
			// FIXME: figure out what to do with weird metadata.
			if strings.HasPrefix(f.Name, "__MACOSX") {
				continue
			}
			
			r, err := f.Open()
			
			blob := &memBlob{
				name: f.Name,
				modtime: f.ModTime(),
				mode: f.Mode(),
				size: int64(f.UncompressedSize64),
				contents: r,
				err: err,
				source: source,
			}
			
			if err != nil {
				result <- blob
				return
			}
			
			blobs = append(blobs, blob)
		}
		
		sort.Sort(blobSorter(blobs))
		
		var prev Blob
		
		for _, blob := range blobs {
			if prev == nil || prev.Name() != blob.Name() {
				// OS X sometimes generates duplicate directories
				result <- blob
			}
			prev = blob
		}
	})
}

// Writes this tree to the given zip file, returning an error on failure.
func (t Tree) ToZip(w *zip.Writer) error {
	for blob := range t {
		fmt.Printf("writing %s\n", blob.Name())
		hdr, err := zip.FileInfoHeader(blob)
		if err != nil {
			return err
		}

		writer, err := w.CreateHeader(hdr)
		if err != nil {
			return err
		}

		if blob.Contents() != nil {
			_, err := io.Copy(writer, blob.Contents())
			if err != nil {
				return err
			}
		}

		if blob.Error() != nil {
			return blob.Error()
		}
	}

	return nil
}