package goscript

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestCompileSimple(t *testing.T) {
	t.Parallel()
	compiler := NewCompiler()
	prog, err := compiler.Compile(CompileJob{
		MainFilePath:       filepath.Join(TESTS, "externals.gs"),
		LocalWorkspaceRoot: TESTS,
		VendorPath:         VENDORPATH,
		StandardLibPath:    STDPATH,
	})
	if err != nil {
		log.Fatal(err)
	}
	f, _ := os.OpenFile("./externals.gob", os.O_RDWR, 0644)
	enc := gob.NewEncoder(f)
	enc.Encode(prog)
	f.Close()
	fmt.Println(prog.String())
}

func TestCompileHelloWorld(t *testing.T) {
	t.Parallel()
	compiler := NewCompiler()
	prog, err := compiler.Compile(CompileJob{
		MainFilePath:       filepath.Join(TESTS, "hello.gs"),
		LocalWorkspaceRoot: TESTS,
		VendorPath:         VENDORPATH,
		StandardLibPath:    STDPATH,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prog.String())
	rt := NewRuntime()
	val := rt.Exec(*prog)
	conv := *val.(*BinaryTypedValue).Value.(*string)
	if conv != "Hello World" {
		t.Fatalf("failed to run hello world, expected 'Hello World' but got '%v'", conv)
	}
}
