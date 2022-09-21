package goscript

import (
	"path/filepath"
	"testing"
)

func TestGenerateFQSC(t *testing.T) {

	ret, err := discoverSources(filepath.Join(TESTS, "externals.gs"), TESTS)
	if err != nil {
		t.Fatalf("sourcewalk failed with error %v", err)
	}
	_, err = generateFQSC(ret)
	if err != nil {
		t.Fatalf("fqsc generation failed with error %v", err)
	}
}
