package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var output string
	if outputFile == "" {
		output = compiler.Compile(string(fileContents))
		fmt.Println(output)
	} else {
		err = ioutil.WriteFile(outputFile, []byte(compiler.Compile(string(fileContents))), 0644)
		if err != nil {
			panic(err)
		}
	}
}
