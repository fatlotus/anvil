package anvil

import (
	"io"
	"os"
)

// A Blob represents a creation or update to a specific path on a volume.
// In particular, if Contents() is non-nil, the file is meant to be opened (or
// created) and then the given contents written to it. Otherwise, the file is
// meant to be erased. If the blob IsDir(), then the blob is a directory instead
// of a traditional file.
type Blob interface {
	Contents() io.Reader // an io.Reader, or nil, if this Blob is a deletion
	Error() error        // if any errors were encountered after or during this record
	Source() string      // where this blob came from
	os.FileInfo
}

// A tree represents a stream of Blobs sorted in lexicographic order.
// To read through the tree, simply iterate through the channel until it it
// closes or one of the entries has the Error() value set. In that case, report
// an error and terminate; continuing to read has undefined behavior.
type Tree <-chan Blob
