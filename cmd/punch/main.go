package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/dfirebaugh/punch/emitters/js"
	"github.com/dfirebaugh/punch/lexer"
	"github.com/dfirebaugh/punch/parser"
	"github.com/sirupsen/logrus"
)

func main() {
	var outputFile string
	var outputTokens bool
	var outputJS bool
	var outputAst bool
	var showHelp bool
	var logLevel string

	flag.StringVar(&outputFile, "o", "", "output file (default: <input_filename>.wasm)")
	flag.BoolVar(&outputTokens, "tokens", false, "output tokens")
	flag.BoolVar(&outputAst, "ast", false, "output Abstract Syntax Tree (AST) file")
	flag.BoolVar(&outputJS, "js", false, "outputs js to stdout")
	flag.BoolVar(&showHelp, "help", false, "show help message")
	flag.StringVar(&logLevel, "log", "error", "set log level (options: trace, debug, info, warn, error, fatal, panic)")
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
	fmt.Println("Usage:", os.Args[0], "[-o output_file] [--tokens] [--wat] [--ast] [--js] [--log log_level] <filename>")
	fmt.Println("Options:")
	fmt.Println("  -o string")
	fmt.Println("        output file (default: <input_filename>.wasm)")
	fmt.Println("  --tokens")
	fmt.Println("        output tokens")
	fmt.Println("  --ast")
	fmt.Println("        output Abstract Syntax Tree (AST) file")
	fmt.Println("  --js")
	fmt.Println("        output Javascript to stdout")
	fmt.Println("  --log string")
	fmt.Println("        set log level (options: trace, debug, info, warn, error, fatal, panic)")
	fmt.Println("  --help")
	fmt.Println("        show help message")
}
