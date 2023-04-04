package repl

import (
	"fmt"
	"io"
	"strings"

	"github.com/dfirebaugh/punch/internal/lexer"
	"github.com/dfirebaugh/punch/internal/parser"
	"github.com/dfirebaugh/punch/internal/wat"

	"github.com/chzyer/readline"
)

const PROMPT = ">> "

type REPL struct {
	in         io.Reader
	out        io.Writer
	cmdHistory []string
}

func New(in io.Reader, out io.Writer) *REPL {
	return &REPL{
		in:  in,
		out: out,
	}
}

func (repl *REPL) Start() {
	rl, err := readline.New(PROMPT)
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	fmt.Println("type 'exit' to close")
	for {
		line, err := rl.Readline()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintln(repl.out, err)
			continue
		}

		line = strings.TrimSpace(line)
		repl.cmdHistory = append(repl.cmdHistory, line)
		shouldContinue := repl.handleLine(line)
		if !shouldContinue {
			break
		}
	}
}

func (repl *REPL) handleLine(line string) bool {
	line = strings.TrimSpace(line)
	switch line {
	case "history":
		repl.showHistory()
	case "clear":
		repl.clearScreen()
	case "exit":
		return false
	default:
		fmt.Fprintf(repl.out, "Command entered: %s\n", line)
		l := lexer.New(line)
		// println("tokens:")
		// for _, tok := range l.Run() {
		// 	println(tok.String())
		// }
		p := parser.New(l)
		program := p.ParseProgram()

		println("ast:")
		json, err := program.JSONPretty()
		if err != nil {
			println(err.Error())
		}
		println(json)

		println("")
		println("wat:")
		println(wat.GenerateWAT(program))
	}
	return true
}

func (repl *REPL) showHistory() {
	for i, cmd := range repl.cmdHistory {
		fmt.Fprintf(repl.out, "%d: %s\n", i+1, cmd)
	}
}

func (repl *REPL) clearScreen() {
	fmt.Fprint(repl.out, "\033[2J\033[H")
}
