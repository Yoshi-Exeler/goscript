package goscript

import (
	"testing"
)

func TestGenerateFQSC(t *testing.T) {
	t.Parallel()
	ret, err := discoverSources("../../tests/externals.gs", "../../tests/")
	if err != nil {
		t.Fatalf("sourcewalk failed with error %v", err)
	}
	_, err = generateFQSC(ret)
	if err != nil {
		t.Fatalf("fqsc generation failed with error %v", err)
	}
}
