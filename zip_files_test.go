package anvil

import (
	"archive/zip"
	"bytes"
	"testing"
)

func TestStreamToZipAndBack(t *testing.T) {
	// write fixture to tree
	buf := new(bytes.Buffer)
	fixtureValidStream().ToZip(zip.NewWriter(buf))

	// read fixture back
	rdr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}

	compareStreams(fixtureValidStream(), FromZip(rdr, "Archive.zip"), t)
}
