package gscompiler

import (
	"testing"
)

func TestSourceWalk(t *testing.T) {
	_, err := DiscoverSources("../../tests/externals.gs", "../../tests/")
	if err != nil {
		t.Fatalf("got error %v", err)
	}
}
