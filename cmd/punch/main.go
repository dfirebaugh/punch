package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dfirebaugh/punch/internal/compiler"
)

func main() {
	var outputFile string
	flag.StringVar(&outputFile, "o", "", "output file")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage:", os.Args[0], "[-o output_file] <filename>")
		os.Exit(1)
	}

	filename := flag.Arg(0)
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	wat, wasm := compiler.Compile(string(fileContents))
	if outputFile == "" {
		fmt.Println(wat)
	} else {
		err = os.WriteFile(outputFile+".wat", []byte(wat), 0644)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(outputFile+".wasm", wasm, 0644)
		if err != nil {
			panic(err)
		}
	}
}
