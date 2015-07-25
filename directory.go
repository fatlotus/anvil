package anvil

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func FromDirectory(root string) Tree {
	root = filepath.Clean(root)

	return makeTree(func(result Tree) {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			var contents *os.File

			if path == root {
				return nil
			}

			if err == nil {
				contents, err = os.Open(path)
				defer contents.Close()
			}

			rel, err := filepath.Rel(root, path)
			sz := info.Size()

			if info.IsDir() {
				rel += "/"
				sz = 0
			}

			result <- &memBlob{
				name:     rel,
				modtime:  info.ModTime(),
				size:     sz,
				mode:     info.Mode(),
				contents: contents,
				err:      err,
				source:   root,
			}

			return err
		})
	}).withValidation()
}

func (t Tree) ToDirectory(root string) error {
	for blob := range t {
		if blob.Error() != nil {
			return blob.Error()
		}

		path := filepath.Join(root, blob.Name())

		switch {
		case blob.Contents() == nil:
			if err := os.RemoveAll(path); err != nil {
				return err
			}
		case blob.IsDir():
			if err := os.MkdirAll(path, blob.Mode()); err != nil {
				return err
			}

			if err := os.Chmod(path, blob.Mode()); err != nil {
				return err
			}

			t := blob.ModTime()

			// Because metadata operations affect the modtime of the parent,
			// we need to defer the call to utimes until the directory is
			// complete; this is a sneaky way to do that.
			defer func() {
				err := os.Chtimes(path, t, t)
				if err != nil {
					panic(err)
				}
			}()

		default:
			if err := os.RemoveAll(path); err != nil {
				return err
			}

			if blob.Mode()&os.ModeSymlink != 0 {

				target, err := ioutil.ReadAll(blob.Contents())
				if err != nil {
					return err
				}

				if err := os.Symlink(string(target), path); err != nil {
					return err
				}

				// TODO: Add support for lchmod and lutimes to Go

			} else {

				fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, blob.Mode())

				if err != nil {
					return err
				}

				err = func() error {
					defer fp.Close()
					if _, err := io.Copy(fp, blob.Contents()); err != nil {
						return err
					}
					return nil
				}()
				if err != nil {
					return err
				}

				if err := os.Chmod(path, blob.Mode()); err != nil {
					return err
				}

				if err := os.Chtimes(path, blob.ModTime(), blob.ModTime()); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
