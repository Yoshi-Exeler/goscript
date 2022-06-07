package goscript

import (
	"testing"
)

// TestGetRequiredExternals parses the externals of the specified file
func TestGetRequiredExternals(t *testing.T) {
	ret, err := getRequiredExternals("../../tests/externals.gs")
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	if len(ret) != 1 {
		t.Fatalf("expected one externals but got %v", len(ret))
	}
}

func TestSourceWalk(t *testing.T) {
	_, err := SourceWalk("../../tests/externals.gs", "../../tests/")
	if err != nil {
		t.Fatalf("got error %v", err)
	}
}
