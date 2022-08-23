package goscript

import (
	"testing"
)

func TestSourceWalk(t *testing.T) {
	_, err := discoverSources("../../tests/externals.gs", "../../tests/")
	if err != nil {
		t.Fatalf("got error %v", err)
	}
}
