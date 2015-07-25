package hephaestus

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// The withValidation transformer ensures that this stream is valid.
func (target Tree) withValidation() Tree {
	return makeTree(func(result Tree) {

		var prev Blob
		var err error

		for blob := range target {
			switch {
			case blob.Error() != nil:
				err = blob.Error()

			case len(blob.Name()) == 0 || blob.Name()[0] == '/':
				err = fmt.Errorf("path %s is invalid\n", blob.Name())

			case prev != nil && !isBlobLessThan(prev, blob):
				err = fmt.Errorf("out of order:\n a:  %s\n vs: %s\n\n",
					formatBlob(prev), formatBlob(blob))

			case blob.IsDir() != strings.HasSuffix(blob.Name(), "/"):
				err = fmt.Errorf("IsDir() = %t != %t = has slash\n",
					blob.IsDir(), strings.HasSuffix(blob.Name(), "/"))
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
func (t Tree) WithPrefix(prefix string) Tree {
	prefix = filepath.Clean(prefix)

	if prefix == "." {
		return t
	}

	return makeTree(func(result Tree) {

		first := true

		for blob := range t {

			// wait until we know Source() before emitting prefix directories.
			if first {
				parts := strings.Split(prefix, "/")

				if parts[0] != "" {
					for i := 0; i < len(parts); i++ {
						path := filepath.Join(parts[0:i+1]...) + "/"
						result <- &memBlob{
							name:     path,
							modtime:  time.Now(),
							contents: strings.NewReader(""),
							source:   "prefix for " + blob.Source(),
							mode:     os.FileMode(0755) | os.ModeDir,
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

// The shallowDiff function determines if two blobs represent the same content.
func shallowDiff(a, b Blob) string {
	dt := a.ModTime().UTC().Unix() - b.ModTime().UTC().Unix()

	if !(-2 <= dt && dt <= 2) {
		return fmt.Sprintf("a.ModTime() = %s != %s = b.ModTime()",
			a.ModTime(), b.ModTime())
	}

	if a.Size() != b.Size() {
		return fmt.Sprintf("a.Size() = %d != %d = b.Size()",
			a.Size(), b.Size())
	}

	if a.Mode() != b.Mode() {
		return fmt.Sprintf("a.Mode() = %s != %s = b.Mode()",
			a.Mode(), b.Mode())
	}

	return ""
}

// The overlayTwo method applies the results of over on top of under.
func overlayTwo(under, over Tree) Tree {
	return makeTree(func(result Tree) {

		hunder := <-under
		hover := <-over

		for hunder != nil || hover != nil {
			if isBlobLessThan(hunder, hover) {
				result <- hunder
				hunder = <-under
			} else if isBlobLessThan(hover, hunder) {
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

// The OverlayAll method applies all subsequent elements on top of the first.
func OverlayAll(trees []Tree) Tree {

	if len(trees) == 1 {
		return trees[0]
	} else if len(trees) == 0 {
		return nil
	}

	p := len(trees) / 2
	return overlayTwo(
		OverlayAll(trees[:p]),
		OverlayAll(trees[p:]),
	)

}

// The Difference method computes changes to get from under to over.
func Difference(under, over Tree) Tree {
	return makeTree(func(result Tree) {

		hunder := <-under
		hover := <-over

		for hunder != nil || hover != nil {
			if isBlobLessThan(hunder, hover) { // removed

				c := copyof(hunder) // mark as deletion
				c.contents = nil
				c.modtime = time.Unix(0, 0)
				c.size = 0
				result <- c
				hunder = <-under

			} else if isBlobLessThan(hover, hunder) { // added

				result <- hover
				hover = <-over

			} else {

				difference := shallowDiff(hover, hunder)

				if difference != "" {
					result <- &memBlob{
						name:     hover.Name(),
						modtime:  hover.ModTime(),
						size:     hover.Size(),
						mode:     hover.Mode(),
						contents: hover.Contents(),
						err:      hover.Error(),
						source:   fmt.Sprintf("%s; %s", hover.Source(), difference),
					}
				}
				hover = <-over
				hunder = <-under

			}
		}

	})
}
