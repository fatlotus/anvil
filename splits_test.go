package anvil

import (
	"testing"
)

func TestSplits(t *testing.T) {
	a, b := fixtureValidTree().Split()
	done := make(chan int)

	go func() {
		compareTrees(a, fixtureValidTree(), t)
		done <- 1
	}()

	go func() {
		compareTrees(b, fixtureValidTree(), t)
		done <- 2
	}()

	<-done
	<-done
}
