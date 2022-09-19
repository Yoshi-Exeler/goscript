package goscript

import (
	"testing"
)

func TestSourceWalk(t *testing.T) {
	t.Parallel()
	_, err := discoverSources("../../tests/externals.gs", "../../tests/")
	if err != nil {
		t.Fatalf("got error %v", err)
	}
}
