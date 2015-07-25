package anvil

import (
	"archive/zip"
	"bytes"
	"testing"
)

func TestTreeToZipAndBack(t *testing.T) {
	// write fixture to tree
	buf := new(bytes.Buffer)
	fixtureValidTree().ToZip(zip.NewWriter(buf))

	// read fixture back
	rdr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}

	compareTrees(fixtureValidTree(), FromZip(rdr, "Archive.zip"), t)
}
