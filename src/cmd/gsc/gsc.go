package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Yoshi-Exeler/goscript/pkg/goscript"
)

func main() {
	vendor := flag.String("vendor", goscript.VENDORPATH, "the path to the vendor directory")
	workspace := flag.String("workspace", "../../tests/", "the path to the root of the workspace")
	standard := flag.String("standard", goscript.STDPATH, "the path to the standard library")
	file := flag.String("file", "", "path to the file to compile")

	flag.Parse()

	fmt.Printf(`Goscript Compiler 0.1`)

	comp := goscript.NewCompiler()

	prog, err := comp.Compile(goscript.CompileJob{
		MainFilePath:       *file,
		VendorPath:         *vendor,
		LocalWorkspaceRoot: *workspace,
		StandardLibPath:    *standard,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(prog.String())

	pb, err := goscript.EncodeProgram(prog)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("out.pb", pb, 0644)
	if err != nil {
		panic(err)
	}

}
