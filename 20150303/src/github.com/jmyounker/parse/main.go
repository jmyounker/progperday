package main

// Parser for a very simple list-like calculator language.  Next version
// will probably incorporate a parser, and then subsequently a compiler
// via llvm.

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var prog = flag.String("p", "", "program")

func main() {
	flag.Parse()
	if *prog == "" {
		log.Print("must supply a program")
		os.Exit(127)
	}

	a, err := parse(newLexer(*prog))
	if err != nil {
		log.Printf("error:", err)
		os.Exit(1)
	}
	fmt.Println("PARSED:")
	fmt.Printf("%s\n", a.String())

	mod, f, err := buildIR(a)
	if err != nil {
		log.Printf("error:", err)
		os.Exit(1)
	}

	fmt.Println("COMPILED:")
	mod.Dump()

	fmt.Println("EXECUTED:")
	res := compileAndRun(mod, f)

	fmt.Printf("%f\n", res)
}
