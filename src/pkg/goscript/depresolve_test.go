package goscript

import (
	"path/filepath"
	"testing"
)

// TestGetRequiredExternals parses the externals of the specified file
func TestGetRequiredExternals(t *testing.T) {
	t.Parallel()
	ret, err := getRequiredExternals(filepath.Join(TESTS, "externals.gs"))
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	if len(ret) != 1 {
		t.Fatalf("expected one externals but got %v", len(ret))
	}
}
