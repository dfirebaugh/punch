package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/dfirebaugh/punch/codegen/js"
	"github.com/dfirebaugh/punch/codegen/lua"
	"github.com/dfirebaugh/punch/codegen/punchgen"
	"github.com/dfirebaugh/punch/lexer"
	"github.com/dfirebaugh/punch/parser"
	"github.com/sirupsen/logrus"
)

func main() {
	var outputFile string
	var outputTokens bool
	var outputAst bool
	var genFormat string
	var showHelp bool
	var logLevel string
	var formatCode bool

	flag.StringVar(&outputFile, "o", "", "output file (default: <input_filename>.wasm)")
	flag.BoolVar(&outputTokens, "tokens", false, "output tokens")
	flag.BoolVar(&outputAst, "ast", false, "output Abstract Syntax Tree (AST) file")
	flag.StringVar(&genFormat, "gen", "", "generate code in specified format (options: js, asm, lua, punch)")
	flag.BoolVar(&showHelp, "help", false, "show help message")
	flag.StringVar(&logLevel, "log", "error", "set log level (options: trace, debug, info, warn, error, fatal, panic)")
	flag.BoolVar(&formatCode, "format", false, "format the code")
	flag.Parse()

	if showHelp {
		printUsage()
		os.Exit(0)
	}

	setLogLevel(logLevel)

	if flag.NArg() < 1 {
		printUsage()
		os.Exit(1)
	}

	filename := flag.Arg(0)
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		logrus.Error(err)
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
	program, err := p.ParseProgram(filename)
	if err != nil {
		logrus.Error(err)
		return
	}
	ast, err := program.JSONPretty()
	if err != nil {
		logrus.Error(err)
	}

	if outputFile == "" {
		outputFile = filename
	}

	if outputAst {
		fmt.Printf("%s\n", ast)
	}

	switch genFormat {
	// case "asm":
	// 	gen := asm.NewCodeGenerator()
	// 	asmCode, err := gen.Generate(program)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 		return
	// 	}
	// 	fmt.Println(asmCode)
	// 	return
	case "js":
		t := js.NewTranspiler()
		jsCode, err := t.Transpile(program)
		if err != nil {
			log.Fatalf("error transpiling to js: %v", err)
		}
		fmt.Printf("%s\n", jsCode)
		return
	case "lua":
		t := lua.NewTranspiler()
		luaCode, err := t.Transpile(program)
		if err != nil {
			log.Fatalf("error transpiling to lua: %v", err)
		}
		fmt.Printf("%s\n", luaCode)
		return
	case "punch":
		t := punchgen.NewGenerator()
		punchCode, err := t.Generate(program)
		if err != nil {
			log.Fatalf("error generating punch code: %v", err)
		}
		fmt.Printf("%s\n", punchCode)
		return
	}

	if formatCode {
		t := punchgen.NewGenerator()
		punchCode, err := t.Generate(program)
		if err != nil {
			log.Fatalf("error generating punch code: %v", err)
		}
		err = os.WriteFile(outputFile, []byte(punchCode), 0o644)
		if err != nil {
			log.Fatalf("error writing punch code to file: %v", err)
		}
		fmt.Printf("Punch code written to %s\n", outputFile)
		return
	}

	// Default to JS if no specific output is requested for now
	t := js.NewTranspiler()
	jsCode, err := t.Transpile(program)
	if err != nil {
		log.Fatalf("error transpiling to js: %v", err)
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
			if _, err := nodeStdin.Write([]byte(jsCode)); err != nil {
				log.Fatalf("error writing js file: %s", err)
			}
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

func setLogLevel(level string) {
	switch level {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.ErrorLevel)
	}
}

func printUsage() {
	fmt.Println("Usage:", os.Args[0], "[-o output_file] [--tokens] [--ast] [--gen format] [--log log_level] <filename>")
	fmt.Println("Options:")
	fmt.Println("  -o string")
	fmt.Println("        output file (default: <input_filename>.wasm)")
	fmt.Println("  --tokens")
	fmt.Println("        output tokens")
	fmt.Println("  --ast")
	fmt.Println("        output Abstract Syntax Tree (AST) file")
	fmt.Println("  --gen string")
	fmt.Println("        generate code in specified format (options: js, asm, lua, punch)")
	fmt.Println("  --log string")
	fmt.Println("        set log level (options: trace, debug, info, warn, error, fatal, panic)")
	fmt.Println("  --help")
	fmt.Println("        show help message")
}
