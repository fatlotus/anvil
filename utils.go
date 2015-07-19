package hephaestus

import (
	"io"
	"os"
	"time"
)

type memBlob struct {
	name     string
	modtime  time.Time
	size     int64
	mode     os.FileMode
	contents io.Reader
	source   string
	err      error
}

func (b *memBlob) Name() string {
	return b.name
}

func (b *memBlob) Size() int64 {
	return b.size
}

func (b *memBlob) Mode() os.FileMode {
	return b.mode
}

func (b *memBlob) ModTime() time.Time {
	return b.modtime
}

func (b *memBlob) IsDir() bool {
	return len(b.name) == 0 || b.name[len(b.name)-1] == '/'
}

func (b *memBlob) Sys() interface{} {
	return nil
}

func (b *memBlob) Contents() io.Reader {
	return b.contents
}

func (b *memBlob) Error() error {
	return b.err
}

func (b *memBlob) Source() string {
	return b.source
}

func copyof(b Blob) *memBlob {
	return &memBlob{
		name:     b.Name(),
		modtime:  b.ModTime(),
		size:     b.Size(),
		mode:     b.Mode(),
		contents: b.Contents(),
		err:      b.Error(),
		source:   b.Source(),
	}
}

func makeTree(handler func(Tree)) Tree {
	result := make(Tree, 0)

	go func() {
		handler(result)
		close(result)
	}()

	return result
}
