package anvil

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestTreeToDirectoryAndBack(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	fixtureValidTree().ToDirectory(tmpdir)
	compareTrees(fixtureValidTree(), FromDirectory(tmpdir), t)
}
