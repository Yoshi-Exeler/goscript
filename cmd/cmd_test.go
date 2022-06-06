package main

import (
	"fmt"
	"goscript/pkg/goscript"
	"testing"
)

// TestGetRequiredExternals parses the externals of the specified file
func TestGetRequiredExternals(t *testing.T) {
	ret, err := goscript.GetRequiredExternals("../tests/externals.gs")
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	if len(ret) != 1 {
		t.Fatalf("expected one external but got %v", len(ret))
	}
	fmt.Printf("%+v\n", *ret[0])
}

func TestSourceWalk(t *testing.T) {
	ret, err := goscript.SourceWalk("../tests/externals.gs", "../tests/")
	if err != nil {
		t.Fatalf("got error %v", err)
	}
	fmt.Printf("%+v\n", ret)
}
