package goscript

import (
	"fmt"
	"testing"
)

func TestGenerateFQSC(t *testing.T) {
	ret, err := SourceWalk("../../tests/externals.gs", "../../tests/")
	if err != nil {
		t.Fatalf("sourcewalk failed with error %v", err)
	}
	source, err := generateFQSC(ret)
	if err != nil {
		t.Fatalf("fqsc generation failed with error %v", err)
	}
	fmt.Printf("%+v\n", source)
}
