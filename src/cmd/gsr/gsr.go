package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Yoshi-Exeler/goscript/src/pkg/goscript"
)

func main() {
	file := flag.String("file", "", "the file to load")

	flag.Parse()

	rt := goscript.NewRuntime()

	f, err := os.OpenFile(*file, os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}

	gob.Register(goscript.BT_ANY)
	gob.Register(goscript.Expression{})
	dec := gob.NewDecoder(f)

	var prog goscript.Program

	err = dec.Decode(&prog)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prog.String())
	v := rt.Exec(prog)

	fmt.Println(v)
}
