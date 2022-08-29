package goscript

import (
	"log"
	"testing"
)

func TestCompileSimple(t *testing.T) {
	compiler := NewCompiler()
	_, err := compiler.Compile(CompileJob{
		MainFilePath:       "../../tests/externals.gs",
		LocalWorkspaceRoot: "../../tests/",
		VendorPath:         VENDORPATH,
		StandardLibPath:    STDPATH,
	})
	if err != nil {
		log.Fatal(err)
	}
}
