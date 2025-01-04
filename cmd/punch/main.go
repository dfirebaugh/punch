package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/dfirebaugh/punch/internal/emitters/js"
	"github.com/dfirebaugh/punch/internal/lexer"
	"github.com/dfirebaugh/punch/internal/parser"
)

func main() {
	var outputFile string
	var outputTokens bool
	var outputJS bool
	var outputAst bool
	var showHelp bool

	flag.StringVar(&outputFile, "o", "", "output file (default: <input_filename>.wasm)")
	flag.BoolVar(&outputTokens, "tokens", false, "output tokens")
	flag.BoolVar(&outputAst, "ast", false, "output Abstract Syntax Tree (AST) file")
	flag.BoolVar(&outputJS, "js", false, "outputs js to stdout")
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

	l := lexer.New(filename, string(fileContents))

	if outputTokens {
		tokens := l.Run()
		for _, token := range tokens {
			fmt.Println(token)
		}
		return
	}

	p := parser.New(l)
	program := p.ParseProgram(filename)
	ast, err := program.JSONPretty()
	if err != nil {
		panic(err)
	}

	if outputFile == "" {
		outputFile = filename
	}

	if outputAst {
		fmt.Printf("%s\n", ast)
	}

	t := js.NewTranspiler()
	jsCode, err := t.Transpile(program)
	if err != nil {
		log.Fatalf("error transpiling to js: %v", err)
	}

	if outputJS {
		fmt.Printf("%s\n", jsCode)
		return
	}

	bunPath, err := exec.LookPath("bun")
	if err != nil {
		log.Printf("bun is not available on the system, trying node: %v", err)
		nodePath, err := exec.LookPath("node")
		if err != nil {
			log.Fatalf("neither bun nor node is available on the system. Please install one of them.")
		}
		cmd := exec.Command(nodePath, "--input-type=module")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		nodeStdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatalf("failed to open pipe to node: %v", err)
		}

		go func() {
			defer nodeStdin.Close()
			nodeStdin.Write([]byte(jsCode))
		}()

		err = cmd.Run()
		if err != nil {
			log.Fatalf("failed to run node: %v", err)
		}
		return
	}

	cmd := exec.Command(bunPath, "-e", jsCode)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatalf("failed to run bun: %v", err)
	}
}

func printUsage() {
	fmt.Println("Usage:", os.Args[0], "[-o output_file] [--tokens] [--wat] [--ast] [--js] <filename>")
	fmt.Println("Options:")
	fmt.Println("  -o string")
	fmt.Println("        output file (default: <input_filename>.wasm)")
	fmt.Println("  --tokens")
	fmt.Println("        output tokens")
	fmt.Println("  --ast")
	fmt.Println("        output Abstract Syntax Tree (AST) file")
	fmt.Println("  --js")
	fmt.Println("        output Javascript to stdout")
	fmt.Println("  --help")
	fmt.Println("        show help message")
}
