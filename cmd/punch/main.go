package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dfirebaugh/punch/internal/compiler"
)

func main() {
	var outputFile string
	var outputWat bool
	var outputAst bool
	var showHelp bool

	flag.StringVar(&outputFile, "o", "", "output file (default: <input_filename>.wasm)")
	flag.BoolVar(&outputWat, "wat", false, "output WebAssembly Text Format (WAT) file")
	flag.BoolVar(&outputAst, "ast", false, "output Abstract Syntax Tree (AST) file")
	flag.BoolVar(&showHelp, "help", false, "show help message")
	flag.Parse()

	if showHelp {
		printUsage()
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		printUsage()
		os.Exit(1)
	}

	filename := flag.Arg(0)
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	wat, wasm, ast := compiler.Compile(filename, string(fileContents))

	if outputFile == "" {
		outputFile = filename
	}

	if outputWat {
		err = os.WriteFile(outputFile+".wat", []byte(wat), 0o644)
		if err != nil {
			panic(err)
		}
	}

	if outputAst {
		err = os.WriteFile(outputFile+".ast", []byte(ast), 0o644)
		if err != nil {
			panic(err)
		}
	}

	err = os.WriteFile(outputFile+".wasm", wasm, 0o644)
	if err != nil {
		panic(err)
	}
}

func printUsage() {
	fmt.Println("Usage:", os.Args[0], "[-o output_file] [--wat] [--ast] <filename>")
	fmt.Println("Options:")
	fmt.Println("  -o string")
	fmt.Println("        output file (default: <input_filename>.wasm)")
	fmt.Println("  --wat")
	fmt.Println("        output WebAssembly Text Format (WAT) file")
	fmt.Println("  --ast")
	fmt.Println("        output Abstract Syntax Tree (AST) file")
	fmt.Println("  --help")
	fmt.Println("        show help message")
}
