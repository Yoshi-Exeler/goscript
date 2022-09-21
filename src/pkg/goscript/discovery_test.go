package goscript

import (
	"path/filepath"
	"testing"
)

func TestSourceWalk(t *testing.T) {

	_, err := discoverSources(filepath.Join(TESTS, "externals.gs"), TESTS)
	if err != nil {
		t.Fatalf("got error %v", err)
	}
}
