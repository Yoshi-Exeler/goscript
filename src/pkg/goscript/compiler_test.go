package goscript

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"
)

func TestCompileVariableIdentity(t *testing.T) {
	compiler := NewCompiler()
	prog, err := compiler.Compile(CompileJob{
		MainFilePath:       filepath.Join(TESTS, "var_identity.gs"),
		LocalWorkspaceRoot: TESTS,
		VendorPath:         VENDORPATH,
		StandardLibPath:    STDPATH,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prog)
	rt := NewRuntime()
	rt.Exec(*prog)
}

func TestCompileImports(t *testing.T) {
	compiler := NewCompiler()
	prog, err := compiler.Compile(CompileJob{
		MainFilePath:       filepath.Join(TESTS, "imports.gs"),
		LocalWorkspaceRoot: TESTS,
		VendorPath:         VENDORPATH,
		StandardLibPath:    STDPATH,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prog)
	rt := NewRuntime()
	rt.Exec(*prog)
}

// func TestCompileCall(t *testing.T) {
// 	compiler := NewCompiler()
// 	prog, err := compiler.Compile(CompileJob{
// 		MainFilePath:       filepath.Join(TESTS, "call.gs"),
// 		LocalWorkspaceRoot: TESTS,
// 		VendorPath:         VENDORPATH,
// 		StandardLibPath:    STDPATH,
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(prog)
// 	rt := NewRuntime()
// 	rt.Exec(*prog)
// }

// func TestCompileArray(t *testing.T) {
// 	compiler := NewCompiler()
// 	_, err := compiler.Compile(CompileJob{
// 		MainFilePath:       filepath.Join(TESTS, "array.gs"),
// 		LocalWorkspaceRoot: TESTS,
// 		VendorPath:         VENDORPATH,
// 		StandardLibPath:    STDPATH,
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func TestCompileHello(t *testing.T) {
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
	fmt.Println(prog)
	// rt := NewRuntime()
	// rt.Exec(*prog)
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
