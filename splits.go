package anvil

import (
	"bytes"
	"io"
)

// Creates two trees from a single Stream.
//
// Note: this function buffers the contents of every
func (t Stream) Split() (Stream, Stream) {
	a := make(chan Blob)
	b := make(chan Blob)

	go func() {
		for blob := range t {

			buf := new(bytes.Buffer)
			content := blob.Contents()

			io.Copy(buf, content)

			a <- &memBlob{
				name:     blob.Name(),
				modtime:  blob.ModTime(),
				size:     blob.Size(),
				mode:     blob.Mode(),
				contents: bytes.NewReader(buf.Bytes()),
				err:      blob.Error(),
				source:   blob.Source(),
			}

			b <- &memBlob{
				name:     blob.Name(),
				modtime:  blob.ModTime(),
				size:     blob.Size(),
				mode:     blob.Mode(),
				contents: bytes.NewReader(buf.Bytes()),
				err:      blob.Error(),
				source:   blob.Source(),
			}

		}

		close(a)
		close(b)
	}()

	return a, b
}
