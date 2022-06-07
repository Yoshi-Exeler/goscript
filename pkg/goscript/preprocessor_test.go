package goscript

import (
	"testing"
)

func TestGenerateFQSC(t *testing.T) {
	ret, err := DiscoverSources("../../tests/externals.gs", "../../tests/")
	if err != nil {
		t.Fatalf("sourcewalk failed with error %v", err)
	}
	_, err = generateFQSC(ret)
	if err != nil {
		t.Fatalf("fqsc generation failed with error %v", err)
	}
}
