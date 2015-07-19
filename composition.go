package hephaestus

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// Running cleanPath returns a relative version of the string without . or ..
func cleanPath(path string) string {
	result := filepath.Clean("/" + path)
	if path != "" && path[len(path)-1] == '/' {
		result += "/"
	}
	return result[1:]
}

// The withValidation transformer ensures that this stream is valid.
func withValidation(target Tree) Tree {
	return makeTree(func(result Tree) {

		var prev Blob
		var err error

		for blob := range target {
			switch {
			case blob.Error() != nil:
				err = blob.Error()
			case len(blob.Name()) == 0 || blob.Name()[0] == '/':
				err = fmt.Errorf("path %s is invalid\n", blob.Name())
			case prev != nil && blob.Name() <= prev.Name():
				err = fmt.Errorf("out of order:\n a:  %s\n vs: %s\n\n",
					formatBlob(prev), formatBlob(blob))
			case CleanPath(blob.Name()) != blob.Name():
				err = fmt.Errorf("unclean path: was %s, should be %s\n",
					blob.Name(), cleanPath(blob.Name()))
			}

			if err != nil {
				copy := copyof(blob)
				copy.err = err
				result <- copy
				break
			}

			result <- blob
			prev = blob
		}
	})
}

// The AddPrefix function transforms a Tree to exist within the given prefix.
func AddPrefix(t Tree, prefix string) Tree {
	return makeTree(func(result Tree) {

		first := true

		for blob := range t {

			// wait until we know Source() before emitting prefix directories.
			if first {
				parts := strings.Split(CleanPath(prefix), "/")
				if parts[0] != "" {
					for i := 0; i < len(parts); i++ {
						path := filepath.Join(parts[0:i+1]...) + "/"
						result <- &memBlob{
							name:     path,
							modtime:  time.Now(),
							contents: strings.NewReader(""),
							source:   blob.Source(),
						}
					}
				}
				first = false
			}

			b := copyof(blob)
			b.name = filepath.Join(prefix, b.name)
			if blob.IsDir() {
				b.name += "/"
			}

			result <- b
		}

	})
}

// The lessThan comparator compares two (possibly nil) blobs.
func lessThan(a, b Blob) bool {
	if b == nil {
		return (a != nil)
	} else if a == nil {
		return false
	}

	return a.Name() < b.Name()
}

// The shallowDiff function determines if two blobs represent the same content.
func shallowDiff(a, b Blob) bool {
	return (a.ModTime().Unix() != b.ModTime().Unix() ||
		a.Size() != b.Size() ||
		a.Mode() != b.Mode())
}

// The overlayTwo method applies the results of over on top of under.
func overlayTwo(under, over Tree) Tree {
	return makeTree(func(result Tree) {

		hunder := <-under
		hover := <-over

		for hunder != nil || hover != nil {
			if lessThan(hunder, hover) {
				result <- hunder
				hunder = <-under
			} else if lessThan(hover, hunder) {
				result <- hover
				hover = <-over
			} else {
				result <- hover
				hover = <-over
				hunder = <-under
			}
		}

	})
}

// The AddOverlay method applies all elements of the slice in order.
func AddOverlay(trees []Tree) Tree {

	if len(trees) == 1 {
		return trees[0]
	} else if len(trees) == 0 {
		return nil
	}

	p := len(trees) / 2
	return overlayTwo(
		Overlay(trees[:p]),
		Overlay(trees[p:]),
	)

}

// The Difference method computes changes to get from under to over.
func Difference(under, over Tree) Tree {
	return makeTree(func(result Tree) {

		hunder := <-under
		hover := <-over

		for hunder != nil || hover != nil {
			if lessThan(hunder, hover) { // removed

				c := copyof(hunder) // mark as deletion
				c.contents = nil
				c.modtime = time.Unix(0, 0)
				c.size = 0
				result <- c
				hunder = <-under

			} else if lessThan(hover, hunder) { // added

				result <- hover
				hover = <-over

			} else {

				if shallowDiff(hover, hunder) {
					result <- hover
				}
				hover = <-over
				hunder = <-under

			}
		}

	})
}
