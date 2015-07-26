package anvil

import (
	"testing"
)

func TestSplits(t *testing.T) {
	a, b := fixtureValidStream().Split()
	done := make(chan int)

	go func() {
		compareStreams(a, fixtureValidStream(), t)
		done <- 1
	}()

	go func() {
		compareStreams(b, fixtureValidStream(), t)
		done <- 2
	}()

	<-done
	<-done
}
