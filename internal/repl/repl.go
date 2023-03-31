package repl

import (
	"fmt"
	"io"
	"punch/internal/lexer"
	"strings"

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

		for _, t := range l.Run() {
			println(t.String())
		}
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
