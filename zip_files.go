package anvil

import (
	"archive/zip"
	"io"
	"sort"
	"strings"
)

const uint32max = 4294967295

// Reads a Stream from the given zip file stream.
func FromZip(r *zip.Reader, source string) Stream {
	return makeStream(func(result chan<- Blob) {

		blobs := make([]Blob, 0)

		for _, f := range r.File {

			// FIXME: figure out what to do with weird metadata.
			if strings.HasPrefix(f.Name, "__MACOSX") {
				continue
			}

			r, err := f.Open()

			// FIXME: deal with bugs in Zip format
			// t := f.ModTime().UTC()
			// utcd := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local).Local()

			sz := int64(f.UncompressedSize64)

			if strings.HasSuffix(f.Name, "/") {
				sz = 0
			}

			blob := &memBlob{
				name:     f.Name,
				modtime:  f.ModTime(),
				mode:     f.Mode(),
				size:     sz,
				contents: r,
				err:      err,
				source:   source,
			}

			if err != nil {
				result <- blob
				return
			}

			blobs = append(blobs, blob)
		}

		sort.Sort(sortBlobs(blobs))

		var prev Blob

		for _, blob := range blobs {
			if prev == nil || prev.Name() != blob.Name() {
				// OS X sometimes generates duplicate directories
				result <- blob
			}
			prev = blob
		}

	}).withValidation()
}

// Writes this tree to the given zip file, returning an error on failure.
func (t Stream) ToZip(w *zip.Writer) error {
	for blob := range t {

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

	w.Close()

	return nil
}
