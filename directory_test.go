package anvil

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestStreamToDirectoryAndBack(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	fixtureValidStream().ToDirectory(tmpdir)
	compareStreams(fixtureValidStream(), FromDirectory(tmpdir), t)
}
