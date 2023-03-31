package compiler

import (
	"punch/internal/lexer"
)

type Compiler struct {
	Source string
}

func New(source string) Compiler {
	return Compiler{
		Source: source,
	}
}

func (c Compiler) Run() {
	c.compile(c.getSourceCode())
}
func (c Compiler) extractCommandLineArguments() {}

func (c Compiler) compile(source string) {
	l := lexer.New(source)
	tokens := l.Run()
	for _, t := range tokens {
		println(t.String())
	}
}

func (c Compiler) getSourceCode() string {
	return c.Source
}
func (c Compiler) reportSyntaxError() {}
func (c Compiler) optimize()          {}
func (c Compiler) generateCode()      {}
func (c Compiler) createGenerator()   {}
