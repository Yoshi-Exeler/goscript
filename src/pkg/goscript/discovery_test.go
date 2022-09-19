package goscript

import (
	"path/filepath"
	"testing"
)

func TestSourceWalk(t *testing.T) {
	t.Parallel()
	_, err := discoverSources(filepath.Join(TESTS, "externals.gs"), TESTS)
	if err != nil {
		t.Fatalf("got error %v", err)
	}
}
