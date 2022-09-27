package goscript

import (
	"log"
	"path/filepath"
	"testing"
)

func TestCompileImports(t *testing.T) {
	compiler := NewCompiler()
	_, err := compiler.Compile(CompileJob{
		MainFilePath:       filepath.Join(TESTS, "imports.gs"),
		LocalWorkspaceRoot: TESTS,
		VendorPath:         VENDORPATH,
		StandardLibPath:    STDPATH,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestCompileArray(t *testing.T) {
	compiler := NewCompiler()
	_, err := compiler.Compile(CompileJob{
		MainFilePath:       filepath.Join(TESTS, "array.gs"),
		LocalWorkspaceRoot: TESTS,
		VendorPath:         VENDORPATH,
		StandardLibPath:    STDPATH,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestCompileHello(t *testing.T) {
	compiler := NewCompiler()
	_, err := compiler.Compile(CompileJob{
		MainFilePath:       filepath.Join(TESTS, "hello.gs"),
		LocalWorkspaceRoot: TESTS,
		VendorPath:         VENDORPATH,
		StandardLibPath:    STDPATH,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestCompileTypecast(t *testing.T) {
	compiler := NewCompiler()
	_, err := compiler.Compile(CompileJob{
		MainFilePath:       filepath.Join(TESTS, "typecast.gs"),
		LocalWorkspaceRoot: TESTS,
		VendorPath:         VENDORPATH,
		StandardLibPath:    STDPATH,
	})
	if err != nil {
		log.Fatal(err)
	}
}
